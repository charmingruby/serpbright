package brightdata

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func BrighDataResultToSearchResult(result BrightDataSearchResult, html string, task entity.CampaignTask) process_entity.SearchResult {
	return process_entity.SearchResult{
		ID:              "",
		Results:         []process_entity.SearchResultItem{},
		ShoppingResults: []process_entity.ShoppingSearchResultItem{},
		Task:            task,
		SearchUrl:       result.General.SearchType,
		HTMLData:        html,
		CreatedAt:       time.Now(),
	}
}

// func (s *BrightData) filterADs() {}
