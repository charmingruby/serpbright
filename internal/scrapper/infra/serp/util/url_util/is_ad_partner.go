package url_util

import (
	"net/url"

	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
)

func IsAdPartner(urlData *url.URL, domain string) bool {

	_, ok := constant.AdPartners[urlData.Host]

	if !ok {
		_, ok = constant.AdPartners[domain+urlData.Path]
	}

	if !ok {
		_, ok = constant.AdPartners[domain]
	}

	return ok
}
