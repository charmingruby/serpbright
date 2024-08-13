package brightdata

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *BrightData) parseResult(apiData data.BrightDataSearchResult, task entity.CampaignTask) process_entity.SearchResult {
	searchResult := process_entity.SearchResult{}

	parser := parser.NewBrightDataParser(parser.BrightDataParserOptions{
		IncludeHTML:                s.SearchConfig.IncludeHTML,
		SkipRedirectAll:            s.SearchConfig.SkipRedirectAll,
		SkipRedirectCampaigns:      s.SearchConfig.SkipRedirectCampaigns,
		SkipCustomerDomainRedirect: s.SearchConfig.SkipCustomerDomainRedirect,
	})

	topADs, bottomADs := parser.FilterADs(&apiData)
	apiData.BottomAds = bottomADs
	apiData.TopAds = topADs

	// ID
	searchResult.ID = primitive.NewObjectID().String()

	// Task
	searchResult.Task = task

	// Engine
	searchResult.SearchUrl = apiData.Input.OriginalURL

	// HTML
	if s.SearchConfig.IncludeHTML {
		searchResult.HTMLData = apiData.HTML
	}

	// Search Results
	searchResult.Results = parser.ParseSearchResults(task, apiData)

	// Shopping ADs

	// Created At
	searchResult.CreatedAt = time.Now()

	return searchResult
}
