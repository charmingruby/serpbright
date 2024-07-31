package adapter

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

type SerpAdapter interface {
	Search(campaigntask entity.CampaignTask) (entity.RawSearchData, error)
}
