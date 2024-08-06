package process_entity

import (
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
)

type SearchResult struct {
	ID        string              `json:"id" bson:"_id"`
	Task      entity.CampaignTask `json:"task" bson:"task"`
	HTMLData  string              `json:"htmlData" bson:"htmlData"`
	SearchUrl string              `json:"searchUrl" bson:"searchUrl"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}
