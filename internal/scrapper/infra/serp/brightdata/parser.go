package brightdata

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/parser"
)

func (s *BrightData) parseResult(apiResult data.BrightDataSearchResult, task entity.CampaignTask) process_entity.SearchResult {
	parser := parser.NewBrightDataParser(s.IncludeHTML)

	topADs, bottomADs := parser.FilterADs(&apiResult)
	apiResult.BottomAds = bottomADs
	apiResult.TopAds = topADs

	searchResult := process_entity.SearchResult{}

	// Engine
	searchResult.SearchUrl = apiResult.Input.OriginalURL

	// HTML
	if s.IncludeHTML {
		searchResult.HTMLData = apiResult.HTML
	}

	// Search Results

	// Shopping ADs

	return process_entity.SearchResult{
		ID:              "",
		Results:         []process_entity.SearchResultItem{},
		ShoppingResults: []process_entity.ShoppingSearchResultItem{},
		Task:            task,
		SearchUrl:       apiResult.General.SearchType,
		HTMLData:        apiResult.HTML,
		CreatedAt:       time.Now(),
	}
}
