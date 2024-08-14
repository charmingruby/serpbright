package parser

import (
	"os"
	"regexp"
	"strings"
)

var intentRegex regexp.Regexp
var gclidRegex regexp.Regexp
var protocolSchemeErrorRegex regexp.Regexp
var blackListedRedirectRegex regexp.Regexp
var skipRedirectRegex regexp.Regexp
var mustRedirectRegex regexp.Regexp

func NewBrightDataParser(opts BrightDataParserOptions) BrightDataParser {
	initRegexes(opts)

	return BrightDataParser{
		SearchOptions: opts,
	}
}

type BrightDataParser struct {
	SearchOptions BrightDataParserOptions
}

type BrightDataParserOptions struct {
	IncludeHTML                bool
	SkipRedirectAll            bool
	SkipRedirectCampaigns      string
	SkipCustomerDomainRedirect bool
	RedirectTimeout            int
	ConcatDomainLastURL        string
	ConcatFirstDomainURL       string
}

func initRegexes(opts BrightDataParserOptions) {
	regex, _ := regexp.Compile(`(?i)^intent:`)
	intentRegex = *regex
	regex, _ = regexp.Compile(`(?i)^gclid=`)
	gclidRegex = *regex
	regex, _ = regexp.Compile(`unsupported protocol scheme`)
	protocolSchemeErrorRegex = *regex
	regex, _ = regexp.Compile(strings.ReplaceAll(regexp.QuoteMeta(os.Getenv("BLACKLISTED_REDIRECT_DOMAINS_REGEX")), "\\|", "|"))
	blackListedRedirectRegex = *regex
	regex, _ = regexp.Compile(`(?i)^(` + opts.SkipRedirectCampaigns + `);`)
	skipRedirectRegex = *regex
	regex, _ = regexp.Compile(`(?i)/aclk?`)
	mustRedirectRegex = *regex
}
