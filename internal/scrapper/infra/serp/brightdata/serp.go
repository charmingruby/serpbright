package brightdata

import (
	"github.com/charmingruby/serpright/internal/common/helper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *BrightData) Search(campaigntask entity.CampaignTask) (process_entity.SearchResult, error) {
	serchResult, err := s.doRequest(campaigntask)
	if err != nil {
		return process_entity.SearchResult{}, err
	}

	if s.DebugMode {
		if err := helper.DebugJSON(serchResult); err != nil {
			return process_entity.SearchResult{}, err
		}
	}

	rawData := BrighDataResultToSearchResult(serchResult)

	return rawData, nil
}
