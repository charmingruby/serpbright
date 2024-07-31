package usecase

import (
	"errors"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
)

func (s *ScrapperUseCaseRegistry) ProcessSerpSearchUseCase(input dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error) {
	campaignTask, err := s.CampaignTaskRepository.FindByID(input.CampaignTaskID)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	rawData, err := s.Serp.Search(*campaignTask)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp search error: " + err.Error())
	}

	processor := entity.NewSearchProcessor(rawData)

	result, err := processor.ProcessData()
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp data process error: " + err.Error())
	}

	return dto.ProcessSerpSearchOutputDTO{
		SearchResult: result,
	}, nil
}
