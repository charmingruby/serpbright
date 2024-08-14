package brightdata

import (
	"log/slog"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/request"
)

func (s *BrightData) Search(campaignTask entity.CampaignTask) (process_entity.SearchResult, error) {
	reqURL := request.BuildBrightDataRequestURL(campaignTask)
	if s.DebugMode {
		slog.Info("BUILT REQUEST URL: " + reqURL)
	}

	searchResult, err := request.DoRequest(reqURL, s.ProxyURL, s.DebugMode)
	if err != nil {
		return process_entity.SearchResult{}, err
	}
	slog.Info("BRIGHT DATA: Processed JSON request")

	rawData := s.parseResult(searchResult, campaignTask)

	return rawData, nil
}
