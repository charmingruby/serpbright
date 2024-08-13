package url_util

import (
	"strings"

	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/util/constant"
)

func ExtractDomain(hostname string) string {
	hostnameParts := strings.Split(hostname, ".")
	hostnamePartsIdx := len(hostnameParts) - 1
	domain := ""
	var hasCCTLD, hasTLD bool = false, false
	for i := hostnamePartsIdx; i >= 0; i-- {
		part := hostnameParts[i]

		if hasCCTLD && hasTLD {
			domain = part + domain
			break
		}

		if i == hostnamePartsIdx {
			if _, ok := constant.CCTLD[part]; ok {
				domain = "." + part
				hasCCTLD = true
				continue
			}
		}

		if _, ok := constant.TLDList[part]; ok {
			domain = "." + part + domain
			hasCCTLD = true
			hasTLD = true
		}

		if hasCCTLD && !hasTLD {
			domain = part + domain
			break
		}
	}
	return domain
}
