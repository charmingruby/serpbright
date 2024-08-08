package brightdata

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func BrighDataResultToSearchResult(result BrightDataSearchResult, html string, task entity.CampaignTask) process_entity.SearchResult {
	return process_entity.SearchResult{
		Task:      task,
		SearchUrl: result.General.SearchType,
		HTMLData:  html,
	}
}
