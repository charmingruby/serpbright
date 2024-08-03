package usecase

import (
	"errors"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *ScrapperUseCaseRegistry) ProcessSerpSearchUseCase(input dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error) {
	rawData, err := s.Serp.Search(input.CampaignTask)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp search error: " + err.Error())
	}

	processor := process_entity.NewSearchProcessor(rawData)

	result, err := processor.ProcessData()
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp data process error: " + err.Error())
	}

	return dto.ProcessSerpSearchOutputDTO{
		SearchResult: result,
	}, nil
}
