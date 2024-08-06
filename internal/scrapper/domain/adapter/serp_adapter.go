package adapter

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

type SerpAdapter interface {
	Search(campaigntask entity.CampaignTask) (process_entity.SearchResult, error)
}
