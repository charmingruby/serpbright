package parser

func NewBrightDataParser(includeHTML bool) BrightDataParser {
	return BrightDataParser{
		IncludeHTML: includeHTML,
	}
}

type BrightDataParser struct {
	IncludeHTML bool
}
