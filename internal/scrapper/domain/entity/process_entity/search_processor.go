package process_entity

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/google/uuid"
)

func NewSearchProcessor(rawData RawSearchData, task entity.CampaignTask) *SearchProcessor {
	return &SearchProcessor{
		Task:    task,
		RawData: rawData,
	}
}

type SearchProcessor struct {
	Task    entity.CampaignTask
	RawData RawSearchData
}

func (sr *SearchProcessor) ProcessData() (SearchResult, error) {
	return SearchResult{
		ID:        uuid.NewString(),
		Task:      sr.Task,
		HTMLData:  "<span>test</span>",
		SearchUrl: sr.RawData.SearchType,
		CreatedAt: time.Now(),
	}, nil
}
