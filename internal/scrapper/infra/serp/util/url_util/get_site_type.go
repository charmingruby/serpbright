package url_util

import "github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"

func GetSiteType(hostname string, domain string, path string) string {
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
