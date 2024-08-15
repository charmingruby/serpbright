package parser

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/charmingruby/serpright/internal/common/helper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/request"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/url_util"
)

func (p *BrightDataParser) ParseSearchResults(
	task entity.CampaignTask,
	apiData data.BrightDataSearchResult,
	ads []ADFilter,
) []process_entity.SearchResultItem {
	results := []process_entity.SearchResultItem{}
	r, _ := regexp.Compile(`(?i)\b` + task.BrandName + `\b`)

	for _, ad := range ads {
		slog.Info(fmt.Sprintf("Redirecting to %s", ad.AD.Link))
		initialRedirectTime := time.Now()

		adURL := ad.AD.Link
		if len(adURL) == 0 && len(ad.AD.DisplayLink) != 0 {
			adURL = ad.AD.DisplayLink
		}

		urlData, _ := url.Parse(p.percentEncodingUrl(adURL))
		hostName := urlData.Hostname()
		domain := strings.ToLower(url_util.ExtractDomain(hostName))

		if len(domain) == 0 {
			domain = strings.ToLower(url_util.ExtractDomain(ad.AD.Link))
		}

		subDomain := strings.Replace(hostName, domain, "", 1)
		redirectedUrls := []string{}

		if p.SearchOptions.SkipRedirectAll || skipRedirectRegex.MatchString(task.ID) || (len(adURL) > 0 && (blackListedRedirectRegex.MatchString(domain) || (p.SearchOptions.SkipCustomerDomainRedirect && domain == strings.ToLower(task.Domain)))) {
			redirectedUrls = append(redirectedUrls, adURL)
		} else {
			redirectedUrls = p.getRedirected(ad.AD.ReferralLink, initialRedirectTime)
			redirectedUrls = p.checkQSUrls(redirectedUrls)
		}

		finishedRedirectTime := time.Since(initialRedirectTime)
		slog.Info(fmt.Sprintf("Redirect levou %s", finishedRedirectTime))

		httpCode := 200
		lastUrl := redirectedUrls[len(redirectedUrls)-1]

		if matched, _ := regexp.MatchString("^(400|404|204)$", lastUrl); matched {
			httpCode, _ = strconv.Atoi(lastUrl)
			redirectedUrls = redirectedUrls[:len(redirectedUrls)-1]
			lastUrlData, _ := url.Parse(redirectedUrls[len(redirectedUrls)-1])
			lastUrlDomain := url_util.ExtractDomain(lastUrlData.Host)
			if url_util.IsAdPartner(lastUrlData, lastUrlDomain) {
				redirectedUrls = append(redirectedUrls, adURL)
			}
		}

		lastUrl = p.percentEncodingUrl(redirectedUrls[len(redirectedUrls)-1])

		sitetype := "default"
		isOwner := false
		phone := ""
		prefix := strings.Split(lastUrl, ":")[0]

		if _, ok := constant.PhonePrefixes[prefix]; ok {
			phone = strings.ReplaceAll(strings.ReplaceAll(lastUrl, prefix+"://", ""), prefix+":", "")
			lastUrl = ad.AD.Link
			if len(lastUrl) == 0 {
				lastUrl = ad.AD.DisplayLink
			}

			if len(lastUrl) == 0 {
				lastUrl = phone
				sitetype = "phone"
			} else {
				prefix = strings.Split(lastUrl, ":")[0]
			}
		}

		if _, ok := constant.SitePrefixes[prefix]; ok {
			lastUrl = p.removeGCLID(lastUrl)

			if strings.HasPrefix(lastUrl, "https://www.google.com.br") {
				lastUrl = ad.AD.Link
			}

			urlData, _ := url.Parse(lastUrl)
			hostName := urlData.Hostname()

			domain = strings.ToLower(url_util.ExtractDomain(hostName))
			isOwner = domain == strings.ToLower(task.Domain)
			sitetype = url_util.GetSiteType(hostName, domain, urlData.Path)

			if sitetype == "phone" {
				if domain == "whatsapp.com" || domain == "wa.me" {
					domain = urlData.Host + urlData.Path + "?phone=" + urlData.Query().Get("phone")
				}
			}

			if sitetype == "searchresult" {
				if strings.Split(domain, ".")[0] == "google" {
					domain = domain + urlData.Path + "?q=" + urlData.Query().Get("q")
				}
			}
			if sitetype == "bizprofile" || sitetype == "sitebuilder" {

				if use, ok := constant.SiteBuilderDomainWithSubDomainAndPath[domain]; ok && use {
					domain = urlData.Host
					if len(urlData.Path) > 1 {
						domain += urlData.Path
					}
				}

				if use, ok := constant.SiteBuilderDomainWithSubDomain[domain]; ok && use {
					domain = urlData.Host
				}
			}

			if sitetype == "appstore" {
				if urlData.Host == "play.google.com" {
					domain = urlData.Host + urlData.Path + "?id=" + urlData.Query().Get("id")
				}
				if urlData.Host == "apps.apple.com" {
					domain = urlData.Host + urlData.Path
				}
			}

			if sitetype == "socialnetwork" {
				if domain == "linktr.ee" {
					domain = domain + urlData.Path
				}

				if domain == "instagram.com/accounts/login/" || domain == "instagram.com" {
					urlData, _ := url.Parse(ad.AD.Link)
					hostName := urlData.Hostname()

					domain = strings.ToLower(url_util.ExtractDomain(hostName)) + urlData.Path
					isOwner = domain == strings.ToLower(task.Domain)
				}
			}

			if sitetype == "default" || sitetype == "incomplete" {
				subDomain = strings.Replace(hostName, domain, "", 1)
			}
		}

		adDomain := url_util.ExtractDomain(ad.AD.Link)
		apiSubDomainString := url_util.ExtractSubdomain(adDomain)
		re := regexp.MustCompile(`^(https?://)?(www\.)?`)
		adDomain = re.ReplaceAllString(adDomain, "")

		if url_util.IsAdPartner(urlData, domain) && adDomain != "" {
			domain = adDomain
			subDomain = apiSubDomainString
		}

		displayedURL := fmt.Sprintf("https://%s%s", subDomain, adDomain)

		campaignID, device, deviceType, keyword, page, geoLocation := task.ExtractDataFromID()

		mongoGeneratedID := primitive.NewObjectID().String()

		adData := process_entity.SearchResultItem{
			Id:             &mongoGeneratedID,
			CampaignTaskId: task.ID,
			CampaignId:     campaignID,
			GeoLocation:    geoLocation,
			Device:         device,
			DeviceType:     deviceType,
			Keyword:        keyword,
			Page:           page,
			// Values
			Position:           uint8(ad.AD.Rank),
			PositionOnPage:     ad.blockPosition,
			Title:              ad.AD.Title,
			BrandInTitle:       r.MatchString(ad.AD.Title),
			BrandInDescription: r.MatchString(ad.AD.Description),
			Url:                lastUrl,
			UrlSequence:        redirectedUrls,
			SiteType:           sitetype,
			TrackingUrl:        ad.AD.ReferralLink,
			DisplayedDomain:    adDomain,
			Domain:             domain,
			Phone:              phone,
			IsOwner:            isOwner,
			SubDomain:          subDomain,
			LinkUrl:            ad.AD.Link,
			DisplayedUrl:       displayedURL,
			Description:        ad.AD.Description,
			Type:               "ad",
			RedirectSeconds:    finishedRedirectTime.Seconds(),
			RedirectHTTPCode:   httpCode,
			Channel:            "search",
			Evidence:           task.HtmlDataUrl,
			CreatedAt:          time.Now()}

		var incompleteFirstUrlDomain = strings.Split(p.SearchOptions.ConcatFirstDomainURL, ",")
		var incompleteDomainFinalUrl = strings.Split(p.SearchOptions.ConcatDomainLastURL, ",")

		p.checkIncompleteFirstUrlDomain(&adData, incompleteFirstUrlDomain)
		p.checkIncompleteDomainFinalUrl(&adData, adData.LinkUrl, incompleteDomainFinalUrl)

		// Coletar API Call de dominios sem tracklink
		base64str := base64.StdEncoding.EncodeToString([]byte(task.GeoLocation))
		if ad.AD.ReferralLink == "" || adDomain == "" {
			uule := url.QueryEscape("w+CAIQICI" + constant.UULEKeys[len(task.GeoLocation)] + base64str)
			google_domain := helper.EmptyString(task.SearchEngineDomain, "google.com.br")
			gl := helper.EmptyString(task.LocaleCountry, "br")
			hl := helper.EmptyString(task.Locale, "pt-br")
			q := url.QueryEscape(task.Keyword)
			page := int(task.Page * request.ItemsPerPage)

			api_call := fmt.Sprintf("&uule=%s&google_domain=%s&gl=%s&hl=%s&q=%s&device=%s&mobile_type=%s&page=%d&geolocation=%s",
				uule,
				google_domain,
				gl,
				hl,
				q,
				request.ExtractDeviceFromTask(task),
				task.MobileType,
				page,
				task.GeoLocation,
			)
			slog.Warn(api_call)
		}

		results = append(results, adData)

	}

	return results
}
