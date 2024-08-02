package usecase

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/adapter"
	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
)

type ScrapperUseCase interface {
	ProcessSerpSearchUseCase(dto dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error)
}

func NewScrapperUseCaseRegistry(
	serp adapter.SerpAdapter) ScrapperUseCaseRegistry {
	return ScrapperUseCaseRegistry{
		Serp: serp,
	}
}

type ScrapperUseCaseRegistry struct {
	Serp adapter.SerpAdapter
}
