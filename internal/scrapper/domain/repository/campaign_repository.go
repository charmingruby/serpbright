package repository

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

type CampaignTaskRepository interface {
	FindByID(id string) (*entity.CampaignTask, error)
}
