package inmemory

import (
	"errors"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
)

var (
	campaigntaskMock = entity.CampaignTask{
		ID: "id",
	}
)

func NewInMemoryCampaignTaskRepository() InMemoryCampaignTaskRepository {
	return InMemoryCampaignTaskRepository{
		Items: []entity.CampaignTask{campaigntaskMock},
	}
}

type InMemoryCampaignTaskRepository struct {
	Items []entity.CampaignTask
}

func (r *InMemoryCampaignTaskRepository) FindByID(id string) (*entity.CampaignTask, error) {
	for _, c := range r.Items {
		if c.ID == id {
			return &c, nil
		}
	}

	return nil, errors.New("campaigntask not found")
}
