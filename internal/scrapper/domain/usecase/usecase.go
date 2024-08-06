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
	searchResultRepo repository.SearchResultRepository,
) *ScrapperUseCaseRegistry {
	return &ScrapperUseCaseRegistry{
		Serp:                   serp,
		SearchResultRepository: searchResultRepo,
	}
}

type ScrapperUseCaseRegistry struct {
	Serp                   adapter.SerpAdapter
	SearchResultRepository repository.SearchResultRepository
}
