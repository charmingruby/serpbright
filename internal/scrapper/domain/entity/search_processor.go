package entity

func NewSearchProcessor(rawData RawSearchData) *SearchProcessor {
	return &SearchProcessor{
		RawData:       rawData,
		ResultantData: nil,
	}
}

type SearchProcessor struct {
	RawData       RawSearchData
	ResultantData *ResultantData
}

func (sr *SearchProcessor) ProcessData() (ResultantData, error) {
	return ResultantData{}, nil
}
