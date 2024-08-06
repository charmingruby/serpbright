package scrapper

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/adapter"
	"github.com/charmingruby/serpright/internal/scrapper/domain/repository"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func NewService(
	serp adapter.SerpAdapter,
	searchResultRepo repository.SearchResultRepository,
) usecase.ScrapperUseCase {
	return usecase.NewScrapperUseCaseRegistry(serp, searchResultRepo)
}
