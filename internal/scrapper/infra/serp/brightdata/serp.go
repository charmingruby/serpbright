package brightdata

import (
	"log/slog"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *BrightData) Search(campaignTask entity.CampaignTask) (process_entity.SearchResult, error) {
	reqURL := s.buildBrightDataRequestURL(campaignTask)
	if s.DebugMode {
		slog.Info("BUILT REQUEST URL: " + reqURL)
	}

	searchResult, err := s.doRequest(reqURL)
	if err != nil {
		return process_entity.SearchResult{}, err
	}
	slog.Info("BRIGHT DATA: Processed JSON request")

	rawData := s.parseResult(searchResult, campaignTask)

	return rawData, nil
}
