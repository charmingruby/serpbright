package usecase

import (
	"errors"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
)

func (s *ScrapperUseCaseRegistry) ProcessSerpSearchUseCase(input dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error) {
	searchResult, err := s.Serp.Search(input.CampaignTask)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	if err := s.SearchResultRepository.Store(searchResult); err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Search result bundle insertion error: " + err.Error())
	}

	if len(searchResult.Results) > 0 {
		if err := s.SearchResultRepository.StoreManyResultItems(searchResult.Results); err != nil {
			return dto.ProcessSerpSearchOutputDTO{}, errors.New("Search result insertion error: " + err.Error())
		}
	}

	if len(searchResult.ShoppingResults) > 0 {
		if err := s.SearchResultRepository.StoreManyShoppingResultItems(searchResult.ShoppingResults); err != nil {
			return dto.ProcessSerpSearchOutputDTO{}, errors.New("Shopping search result insertion error: " + err.Error())
		}
	}

	return dto.ProcessSerpSearchOutputDTO{
		SearchResult: searchResult,
	}, nil
}
