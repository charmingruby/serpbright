package parser

import (
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/url_util"
)

var spaceRegex *regexp.Regexp = regexp.MustCompile("[\\s\u00A0]")

func (p *BrightDataParser) ParseShoppingSearchResults(
	task entity.CampaignTask,
	apiData data.BrightDataSearchResult,
) []process_entity.ShoppingSearchResultItem {
	r, _ := regexp.Compile(`(?i)\b` + task.BrandName + `\b`)
	results := []process_entity.ShoppingSearchResultItem{}

	for _, shoppingAD := range apiData.TopPla {
		err := func(shopAd data.TopPla) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = r.(error)
				}
			}()
			shopUrl := shopAd.Link

			var redirectedUrls []string

			httpCode := 200
			if mustRedirectRegex.MatchString(shopUrl) {
				redirectedUrls = p.getRedirected(shopUrl, time.Now())
				lastUrl := redirectedUrls[len(redirectedUrls)-1]
				if matched, _ := regexp.MatchString("^(400|404|204)$", lastUrl); matched {
					httpCode, _ = strconv.Atoi(lastUrl)
					redirectedUrls = redirectedUrls[:len(redirectedUrls)-1]
					urlData, _ := url.Parse(p.percentEncodingUrl(redirectedUrls[len(redirectedUrls)-1]))
					domain := url_util.ExtractDomain(urlData.Host)
					if url_util.IsAdPartner(urlData, domain) {
						redirectedUrls = append(redirectedUrls, p.getRedirected(lastUrl, time.Now())...)
					}
				} else {
					redirectedUrls = p.checkQSUrls(redirectedUrls)
				}
			} else {
				redirectedUrls = []string{shopUrl}
			}

			lastUrl := p.removeGCLID(redirectedUrls[len(redirectedUrls)-1])

			if strings.HasPrefix(lastUrl, constant.GoogleDomain) {
				lastUrl = shopAd.Link
			}

			urlData, _ := url.Parse(p.percentEncodingUrl(lastUrl))
			hostName := urlData.Hostname()
			domain := strings.ToLower(url_util.ExtractDomain(hostName))
			subDomain := strings.Replace(hostName, domain, "", 1)
			siteType := url_util.CheckSiteType(lastUrl, domain)

			currency, price := p.parsePrice(shopAd.Price)

			adData := process_entity.ShoppingSearchResultItem{
				CompanyName: shopAd.Shop,
				Currency:    currency,
				Price:       price,
				// TODO: TENTAR SIMULAR REGULAR, PROMOÇAO
				//RegularPriceString: shopAd.RegularPrice,
				//RegularPrice:       priceRP,
				PriceString: shopAd.Price,
				SearchResultItem: process_entity.SearchResultItem{
					Position:         uint8(shopAd.Rank),
					PositionOnPage:   "top",
					Title:            shopAd.Title,
					BrandInTitle:     r.MatchString(shopAd.Title),
					Url:              lastUrl,
					UrlSequence:      redirectedUrls,
					SiteType:         siteType,
					TrackingUrl:      shopUrl,
					Domain:           domain,
					SubDomain:        subDomain,
					Type:             "ad",
					Channel:          "shopping",
					RedirectHTTPCode: httpCode,
					CreatedAt:        time.Now(),
				},
			}

			var incompleteFirstUrlDomain = strings.Split(p.SearchOptions.ConcatFirstDomainURL, ",")
			var incompleteDomainFinalUrl = strings.Split(p.SearchOptions.ConcatDomainLastURL, ",")

			p.checkIncompleteFirstUrlDomain(&adData.SearchResultItem, incompleteFirstUrlDomain)
			p.checkIncompleteDomainFinalUrl(&adData.SearchResultItem, adData.Url, incompleteDomainFinalUrl)

			results = append(results, adData)
			return nil
		}(shoppingAD)

		if err != nil {
			slog.Error(err.Error())
		}
	}

	return results
}

func (p *BrightDataParser) parsePrice(price string) (string, float64) {
	price = spaceRegex.ReplaceAllString(price, " ")
	priceParts := strings.Split(price, " ")
	currency := priceParts[0]
	priceParts = strings.Split(price, currency)
	if len(priceParts) < 2 {
		return currency, 0
	}
	priceString := strings.Trim(priceParts[1], " ")
	thousandsSeparator := ","
	for i := len(priceString) - 1; i >= 0; i-- {
		if priceString[i] == '.' {
			break
		}
		if priceString[i] == ',' {
			thousandsSeparator = "."
			break
		}
	}

	priceString = strings.ReplaceAll(strings.ReplaceAll(priceString, thousandsSeparator, ""), ",", ".")

	priceParsed, _ := strconv.ParseFloat(priceString, 64)
	return currency, priceParsed
}

func (p *BrightDataParser) AddShoppingResultItems(shoppingResults []process_entity.ShoppingSearchResultItem) []process_entity.SearchResultItem {
	domainDict := map[string]process_entity.SearchResultItem{}

	for i := 0; i < len(shoppingResults); i++ {
		sRes := shoppingResults[i]
		currentItem, ok := domainDict[sRes.Domain]
		if !ok || sRes.Position < currentItem.Position {
			domainDict[sRes.Domain] = sRes.SearchResultItem
		}
	}

	result := []process_entity.SearchResultItem{}
	for _, v := range domainDict {
		result = append(result, v)
	}

	return result
}
