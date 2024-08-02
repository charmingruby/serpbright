package scrapper

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/adapter"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func NewService(
	serp adapter.SerpAdapter,
) usecase.ScrapperUseCase {
	return &usecase.ScrapperUseCaseRegistry{
		Serp: serp,
	}
}
