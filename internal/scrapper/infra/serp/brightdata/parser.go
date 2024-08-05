package brightdata

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"

type BrightDataResult struct{}

func BrighdataResultToRawData(bdResult BrightDataResult) process_entity.RawSearchData {
	return process_entity.RawSearchData{}
}
