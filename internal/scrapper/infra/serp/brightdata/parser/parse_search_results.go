package parser

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
)

func (p *BrightDataParser) ParseSearchResults(task entity.CampaignTask, apiData data.BrightDataSearchResult) []process_entity.SearchResultItem {
	results := []process_entity.SearchResultItem{}

	return results
}
