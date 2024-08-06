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

	processor := process_entity.NewSearchProcessor(rawData, input.CampaignTask)

	processedResult, err := processor.ProcessData()
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp data process error: " + err.Error())
	}

	//	if err := s.SearchResultRepository.Store(processedResult); err != nil {
	//		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp result insertion error: " + err.Error())
	//}

	return dto.ProcessSerpSearchOutputDTO{
		SearchResult: processedResult,
	}, nil
}
