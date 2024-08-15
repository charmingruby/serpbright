package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/url_util"
	netHTML "golang.org/x/net/html"
)

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
