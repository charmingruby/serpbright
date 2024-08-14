package url_util

import (
	"net/url"
	"strings"
)

func ExtractSubdomain(domain string) string {
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	parsedURL, err := url.Parse(domain)
	if err != nil {
		return ""
	}

	host := parsedURL.Hostname()

	parts := strings.Split(host, ".")

	if len(parts) < 3 {
		return ""
	}

	subdomain := parts[0] + "."
	return subdomain
}
