package process_entity

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/google/uuid"
)

func NewSearchProcessor(searchResult SearchResult, task entity.CampaignTask) *SearchProcessor {
	return &SearchProcessor{
		Task:         task,
		SearchResult: searchResult,
	}
}

type SearchProcessor struct {
	Task         entity.CampaignTask
	SearchResult SearchResult
}

func (sr *SearchProcessor) ProcessData() (SearchResult, error) {
	return SearchResult{
		ID:        uuid.NewString(),
		Task:      sr.Task,
		HTMLData:  "<span>test</span>",
		SearchUrl: sr.SearchResult.SearchUrl,
		CreatedAt: time.Now(),
	}, nil
}
