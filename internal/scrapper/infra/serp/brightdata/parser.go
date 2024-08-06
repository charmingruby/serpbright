package brightdata

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func BrighDataResultToSearchResult(result BrightDataSearchResult) process_entity.SearchResult {
	return process_entity.SearchResult{
		SearchUrl: result.General.SearchType,
	}
}
