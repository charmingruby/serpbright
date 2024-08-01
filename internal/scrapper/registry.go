package scrapper

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/adapter"
	"github.com/charmingruby/serpright/internal/scrapper/domain/repository"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func NewService(

	serp adapter.SerpAdapter,
	campaingTaskRepo repository.CampaignTaskRepository,
) usecase.ScrapperUseCase {
	return &usecase.ScrapperUseCaseRegistry{
		Serp:                   serp,
		CampaignTaskRepository: campaingTaskRepo,
	}
}
