package parser

func NewBrightDataParser(opts BrightDataParserOptions) BrightDataParser {
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
}
