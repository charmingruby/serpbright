package url_util

import (
	"net/url"
	"strings"

	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
)

func CheckSiteType(lastUrl string, domain string) string {
	sitetype := "default"
	subDomain := ""
	prefix := strings.Split(lastUrl, ":")[0]
	if _, ok := constant.PhonePrefixes[prefix]; ok {
		sitetype = "phone"
		domain = lastUrl
	}

	if _, ok := constant.SitePrefixes[prefix]; ok {
		urlData, _ := url.Parse(lastUrl)
		hostName := urlData.Hostname()

		//isOwner = domain == strings.ToLower(task.Domain)

		sitetype = GetSiteType(hostName, domain, urlData.Path)

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

			if domain == "instagram.com" {
				domain = domain + urlData.Path
			}
		}

		subDomain = strings.Replace(hostName, domain, "", 1)
	}
	_ = subDomain
	return sitetype
}
