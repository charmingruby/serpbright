package parser

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	netHTML "golang.org/x/net/html"

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
			sitetype = p.getSiteType(hostName, domain, urlData.Path)

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

		adData := process_entity.SearchResultItem{
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

func (p *BrightDataParser) checkIncompleteFirstUrlDomain(adData *process_entity.SearchResultItem, incompleteFirstUrlDomain []string) {
	for _, domain := range incompleteFirstUrlDomain {
		if adData.Domain == strings.TrimSpace(domain) {
			adData.Domain = adData.SubDomain + adData.Domain
		}
	}
}

func (p *BrightDataParser) checkIncompleteDomainFinalUrl(adData *process_entity.SearchResultItem, urlToCheck string, incompleteDomainFinalUrl []string) {
	for _, domain := range incompleteDomainFinalUrl {
		if adData.Domain == strings.TrimSpace(domain) {
			parsedUrl, err := url.Parse(urlToCheck)
			if err != nil {
				return
			}
			path := strings.Trim(parsedUrl.Path, "/")
			if path != "" {
				adData.Domain = adData.Domain + "/" + path
				break
			}
		}
	}
}

func (p *BrightDataParser) removeGCLID(lastUrl string) string {
	urlParts := strings.Split(lastUrl, "?")
	if len(urlParts) == 1 {
		urlParts = append(urlParts, "")
	}
	qsParams := strings.Split(urlParts[1], "&")
	qsWithoutGCLID := []string{}
	for i := 0; i < len(qsParams); i++ {
		param := qsParams[i]
		if gclidRegex.MatchString(param) {
			continue
		}
		qsWithoutGCLID = append(qsWithoutGCLID, param)
	}
	qs := strings.Join(qsWithoutGCLID, "&")
	if len(qs) > 0 {
		qs = "?" + qs
	}
	return strings.Join([]string{urlParts[0], qs}, "")
}

func (p *BrightDataParser) getSiteType(hostname string, domain string, path string) string {
	if _, ok := constant.SearchEngines[hostname]; ok {
		return "searchengine"
	}

	if _, ok := constant.SearchEngines[domain]; ok {
		return "searchengine"
	}

	if _, ok := constant.FraudSites[hostname]; ok {
		return "fraud"
	}

	if _, ok := constant.FraudSites[domain]; ok {
		return "fraud"
	}

	if v, ok := constant.SiteTypes[hostname+path]; ok {
		return v
	}

	if v, ok := constant.SiteTypes[domain+path]; ok {
		return v
	}

	if v, ok := constant.SiteTypes[hostname]; ok {
		return v
	}

	if v, ok := constant.SiteTypes[domain]; ok {
		return v
	}
	return "default"
}

func (p *BrightDataParser) percentEncodingUrl(urlStr string) string {
	percentEncoded := ""
	urlAndQs := strings.Split(urlStr, "?")
	urlStr = urlAndQs[0]
	urlStrTmp, err := url.PathUnescape(urlStr)
	if err == nil {
		urlStr = urlStrTmp
	}

	urlParts := strings.Split(urlStr, "/")
	for i := 0; i < len(urlParts); i++ {
		if len(percentEncoded) > 0 {
			percentEncoded += "/"
		}
		percentEncoded += url.PathEscape(urlParts[i])
	}
	for i := 1; i < len(urlAndQs); i++ {
		if len(urlAndQs[i]) > 0 {
			percentEncoded += "?" + urlAndQs[i]
		}
	}

	return percentEncoded
}

func (p *BrightDataParser) redirectTimeout() time.Duration {
	timeout := p.SearchOptions.RedirectTimeout

	if timeout <= 0 {
		return time.Duration(10) * time.Second
	}

	return time.Duration(timeout) * time.Second
}

func (p *BrightDataParser) getRedirected(trackingUrl string, sw time.Time) []string {
	redirectTimeout := p.redirectTimeout()

	if time.Since(sw) > redirectTimeout {
		return []string{}
	}

	googlePlayUrl := IsAndroidIntent(trackingUrl)
	if len(googlePlayUrl) > 0 {
		return []string{googlePlayUrl}
	}
	redirectedUrls := []string{trackingUrl}

	req, err := http.NewRequest("GET", trackingUrl, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	client := http.Client{
		Timeout: redirectTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return errors.New("too many redirects")
			}
			adUrl := req.URL.Query().Get("adurl")
			if len(adUrl) > 0 {
				redirectedUrls = append(redirectedUrls, adUrl)
				return http.ErrUseLastResponse
			}

			redirectedUrls = append(redirectedUrls, strings.Join([]string{req.URL.Scheme, "://", req.URL.Host, req.URL.RequestURI()}, ""))
			domain := url_util.ExtractDomain(req.URL.Host)
			if domain == "doubleclick.net" {
				urlStr := doubleclickUrl(redirectedUrls[len(redirectedUrls)-1])
				if len(urlStr) > 0 {
					redirectedUrls = append(redirectedUrls, urlStr)

					return http.ErrUseLastResponse
				}
			}

			return nil
		},
	}

	response, err := client.Do(req)
	lastIdx := len(redirectedUrls) - 1
	if err != nil {
		if protocolSchemeErrorRegex.MatchString(err.Error()) {
			googlePlayUrl := IsAndroidIntent(redirectedUrls[lastIdx])
			if len(googlePlayUrl) > 0 {
				redirectedUrls = append(redirectedUrls, googlePlayUrl)
			}
			return redirectedUrls
		}

		// erros de timeout não podem parar o processo
		if os.IsTimeout(err) {
			return redirectedUrls
		}

		slog.Error(err.Error())
		return redirectedUrls
	}

	if response == nil {
		if strings.HasPrefix(redirectedUrls[lastIdx], "tel:") {
			return redirectedUrls
		}
	}
	if response.StatusCode == http.StatusBadRequest {
		redirectedUrls = append(redirectedUrls, "400")
		return redirectedUrls
	}

	if response.StatusCode == http.StatusNotFound {
		redirectedUrls = append(redirectedUrls, "404")
		return redirectedUrls
	}
	if response.StatusCode == http.StatusNoContent {
		redirectedUrls = append(redirectedUrls, "204")
		return redirectedUrls
	}

	lastUrl := redirectedUrls[lastIdx]
	if matched, _ := regexp.MatchString("{gclid}", lastUrl); matched {
		lastUrl = strings.ReplaceAll(lastUrl, "{gclid}", "")
		redirectedUrls[lastIdx] = lastUrl
	}

	if lastUrl == trackingUrl && response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error(err.Error())
			return redirectedUrls
		}
		urlData, _ := url.Parse(p.percentEncodingUrl(lastUrl))
		domain := url_util.ExtractDomain(urlData.Host)
		ok := url_util.IsAdPartner(urlData, domain)
		urlFromHtml := ""
		htmlStr := string(bodyBytes)
		if ok {
			matched, _ := regexp.MatchString(`^(google\.com|doubleclick\.net)`, domain)
			if matched {
				urlStr := doubleclickUrl(lastUrl)
				if len(urlStr) > 0 {
					redirectedUrls = append(redirectedUrls, urlStr)
					return redirectedUrls
				}
				urlFromHtml = p.getRedirectFromDoubleClickScript(htmlStr)
				if len(urlFromHtml) > 0 {
					query := urlData.Query()
					urlFromQuery := query.Get("h")
					if matched, _ := regexp.MatchString(`^(http|https)://`, urlFromQuery); matched {
						urlFromHtml = urlFromQuery
					}
				}
			} else {
				if domain == "socialsoul.com.vc" {
					urlFromHtml = p.getRedirectFromSocialSoulScript(htmlStr)
				} else {
					newRedirectUrls := p.checkQSUrls(redirectedUrls)
					if len(newRedirectUrls) > lastIdx+1 {
						newRedirectUrls = p.checkBotShield(newRedirectUrls, sw)
						return newRedirectUrls
					}
				}
			}
			if len(strings.Trim(urlFromHtml, " ")) > 0 {
				redirectedUrls = append(redirectedUrls, urlFromHtml)
			}
		}
		lastIdx = len(redirectedUrls) - 1
	}

	googlePlayUrl = IsAndroidIntent(redirectedUrls[lastIdx])
	if len(googlePlayUrl) > 0 {
		redirectedUrls = append(redirectedUrls, googlePlayUrl)
		return redirectedUrls
	}

	urlData, _ := url.Parse(p.percentEncodingUrl(redirectedUrls[lastIdx])) //extractDomain(redirectedUrl)
	domain := url_util.ExtractDomain(urlData.Host)
	ok := url_util.IsAdPartner(urlData, domain)
	if ok {
		lastUrl := redirectedUrls[lastIdx]
		matched, _ := regexp.MatchString(`^(google\.com|doubleclick\.net)`, domain)
		if matched {
			urlStr := doubleclickUrl(lastUrl)
			if len(urlStr) > 0 {
				redirectedUrls = append(redirectedUrls, urlStr)

				return redirectedUrls
			}
		}

		if domain == "clickcease.com" {
			slog.Info("Utilizando clickcease.")
			redirectedUrls = append(redirectedUrls, urlData.Query().Get("url"))
			return redirectedUrls
		}

		newRedirectUrls := p.checkQSUrls(redirectedUrls)
		if len(newRedirectUrls) > lastIdx+1 {
			newRedirectUrls = p.checkBotShield(newRedirectUrls, sw)
			return newRedirectUrls
		}
		var tmpUrls = p.getRedirected(lastUrl, sw)
		if len(tmpUrls) > 1 {
			redirectedUrls = append(redirectedUrls, tmpUrls[1:]...)
		}
	}

	redirectedUrls = p.checkBotShield(redirectedUrls, sw)

	return redirectedUrls
}

func (p *BrightDataParser) checkBotShield(redirectedUrls []string, sw time.Time) []string {
	lastUrl := redirectedUrls[len(redirectedUrls)-1]
	param, found := p.getBotShieldParam(lastUrl)
	if !found {
		return redirectedUrls
	}

	if len(param) == 0 {
		redirectedUrls = redirectedUrls[:len(redirectedUrls)-1]
		return redirectedUrls
	}

	urlData, _ := url.Parse(p.percentEncodingUrl(lastUrl))
	paramValue := urlData.Query().Get(param)

	if len(paramValue) > 0 {
		redirectedUrls = append(redirectedUrls, paramValue)
	} else {
		redirectedUrls = append(redirectedUrls, p.getRedirected(lastUrl, sw)[1:len(redirectedUrls)-1]...)
	}

	return redirectedUrls
}

func (p *BrightDataParser) checkQSUrls(redirectedUrls []string) []string {
	lastUrl := redirectedUrls[len(redirectedUrls)-1]
	urlData, _ := url.Parse(p.percentEncodingUrl(lastUrl))
	q := urlData.Query()
	hostname := urlData.Hostname()
	domain := strings.ToLower(url_util.ExtractDomain(hostname))
	if constant.AdPartners[domain] {
		for i := 0; i < len(constant.AdPartnerParams); i++ {
			val := q.Get(constant.AdPartnerParams[i])
			if len(val) > 0 {
				param, found := p.getBotShieldParam(val)
				if found && len(param) > 0 {
					botShieldParts := strings.Split(lastUrl, val)
					if len(botShieldParts) > 1 {
						val = val + botShieldParts[1]
					}
				}
				redirectedUrls = append(redirectedUrls, val)
				break
			}
		}
	}
	return redirectedUrls
}

func (p *BrightDataParser) getBotShieldParam(lastUrl string) (string, bool) {
	urlData, _ := url.Parse(p.percentEncodingUrl(lastUrl))
	param, ok := constant.BotShield[urlData.Host]
	domain := url_util.ExtractDomain(urlData.Host)
	if !ok {
		param, ok = constant.BotShield[domain]
	}

	return param, ok
}

func (p *BrightDataParser) getRedirectFromDoubleClickScript(html string) string {
	bodyString := "<html><head>" + html + "</head></html>"
	htmlNode, err := netHTML.Parse(strings.NewReader(bodyString))
	if err != nil {
		slog.Error(err.Error())
		return ""
	}

	for {
		if htmlNode == nil {
			return ""
		}

		if htmlNode.Data == "noscript" {
			htmlNode = htmlNode.FirstChild
			break
		}

		if htmlNode.Data != "script" {
			htmlNode = htmlNode.FirstChild
			continue
		}

		if htmlNode.Data == "script" {
			htmlNode = htmlNode.NextSibling
		}
	}

	var metaTag MetaTag
	err = xml.NewDecoder(bytes.NewBuffer([]byte(htmlNode.Data + "</META>"))).Decode(&metaTag)
	if err != nil {
		slog.Error(err.Error())
		return ""
	}

	urlAttr := strings.Split(metaTag.Content, ";")[1]
	urlAttr = strings.Replace(urlAttr, "URL='", "", 1)
	contentUrl, err := url.Parse(p.percentEncodingUrl(urlAttr[:len(urlAttr)-1]))
	qs := contentUrl.Query()
	adurl := qs.Get("adurl")
	if len(adurl) == 0 {
		adurl = qs.Get("ds_dest_url")
	}
	if len(adurl) == 0 {
		adurl = urlAttr
	}
	if err != nil {
		slog.Error(err.Error())
		return ""
	}

	if strings.Contains(adurl, "://ad.doubleclick.net") {
		contentUrl, _ = url.Parse(p.percentEncodingUrl(adurl))
		if matched, _ := regexp.MatchString(`^(http|https)://`, contentUrl.RawQuery); matched {
			adurl = contentUrl.RawQuery
		}
	}

	return adurl
}

func (p *BrightDataParser) getRedirectFromSocialSoulScript(html string) string {
	htmlNode, err := netHTML.Parse(strings.NewReader(html))
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	urlAttr := ""
	for {
		if htmlNode == nil {
			return ""
		}

		if htmlNode.Type == netHTML.DocumentNode && htmlNode.Data == "" {
			htmlNode = htmlNode.LastChild
		}

		if htmlNode.Type == netHTML.ElementNode && htmlNode.Data == "html" || htmlNode.Data == "head" {
			htmlNode = htmlNode.FirstChild
			continue
		}
		if strings.Trim(htmlNode.Data, " ") == "" {
			htmlNode = htmlNode.NextSibling
			continue
		}

		if htmlNode.Data == "meta" {
			for i := 0; i < len(htmlNode.Attr); i++ {
				attr := htmlNode.Attr[i]
				matched, _ := regexp.MatchString("url=", attr.Val)
				if matched {
					urlAttr = attr.Val
					break
				}
			}
			if len(urlAttr) > 0 {
				break
			}
			htmlNode = htmlNode.NextSibling
			continue
		}

		if htmlNode.Type == netHTML.ElementNode || htmlNode.Type == netHTML.CommentNode || htmlNode.Type == netHTML.TextNode {
			htmlNode = htmlNode.NextSibling
		}

	}

	urlAttrparts := strings.Split(urlAttr, ";")
	for i := 0; i < len(urlAttrparts); i++ {
		attrPart := strings.Split(urlAttrparts[i], "url=")
		if len(attrPart) == 2 {
			return attrPart[1]
		}
	}
	return ""
}

func doubleclickUrl(urlStr string) string {
	pathParts := strings.Split(urlStr, ";")
	urlTmp := ""
	if len(pathParts) > 0 {
		urlTmp, _ = url.PathUnescape(pathParts[len(pathParts)-1])
		if len(urlTmp) > 0 {
			tmpParts := strings.Split(urlTmp, "ltd=")
			if tmpParts[0] == "" && len(tmpParts) > 1 {
				urlTmp = strings.Join(tmpParts[1:], "ltd=")
			}
			if len(urlTmp) > 0 && string(urlTmp[0]) != "h" {
				urlTmp = urlTmp[1:]
				if matched, _ := regexp.MatchString(`^(http|https)://`, urlTmp); !matched {
					urlTmp = ""
				}
			}
		}
	}

	if urlStr == urlTmp {
		urlTmp = ""
	}
	return urlTmp
}

func IsAndroidIntent(trackingUrl string) string {
	if !intentRegex.MatchString(trackingUrl) {
		parsedUrl, err := url.Parse(trackingUrl)
		if err != nil {
			slog.Error(err.Error())
			return ""
		}
		trackingUrl = parsedUrl.Query().Get("adurl")
	}

	if !intentRegex.MatchString(trackingUrl) {
		return ""
	}

	parsedUrl, err := url.Parse(strings.Replace(trackingUrl, "intent:", "http:", 1))
	if err != nil {
		slog.Error(err.Error())
		return ""
	}

	return "https://play.google.com/store/apps/details?id=" + parsedUrl.Query().Get("id")
}

type MetaTag struct {
	Text      string `xml:",chardata"`
	HTTPEquiv string `xml:"http-equiv,attr"`
	Content   string `xml:"content,attr"`
}
