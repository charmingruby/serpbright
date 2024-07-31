package usecase

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/adapter"
	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/repository"
)

type ScrapperUseCase interface {
	ProcessSerpSearchUseCase(dto dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error)
}

func NewScrapperUseCaseRegistry(
	serp adapter.SerpAdapter,
	campaigntaskRepo repository.CampaignTaskRepository) ScrapperUseCaseRegistry {
	return ScrapperUseCaseRegistry{
		Serp:                   serp,
		CampaignTaskRepository: campaigntaskRepo,
	}
}

type ScrapperUseCaseRegistry struct {
	Serp                   adapter.SerpAdapter
	CampaignTaskRepository repository.CampaignTaskRepository
}
