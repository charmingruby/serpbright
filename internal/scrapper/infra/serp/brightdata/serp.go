package brightdata

import (
	"fmt"
	"log/slog"

	"github.com/charmingruby/serpright/internal/common/helper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *BrightData) Search(campaignTask entity.CampaignTask) (process_entity.SearchResult, error) {
	reqURL := s.buildBrightDataRequestURL(campaignTask)
	if s.DebugMode {
		slog.Info("BUILT REQUEST URL: " + reqURL)
	}

	var htmlResult string
	if s.IncludeHTML {
		res, err := s.doHTMLRequest(reqURL)
		if err != nil {
			return process_entity.SearchResult{}, err
		}
		htmlResult = res
	}

	if s.DebugMode && s.IncludeHTML {
		fmt.Println(htmlResult)
	}

	serchResult, err := s.doJSONRequest(reqURL)
	if err != nil {
		return process_entity.SearchResult{}, err
	}
	if s.DebugMode {
		if err := helper.DebugJSON(serchResult); err != nil {
			return process_entity.SearchResult{}, err
		}
	}

	rawData := BrighDataResultToSearchResult(serchResult, htmlResult, campaignTask)

	return rawData, nil
}
