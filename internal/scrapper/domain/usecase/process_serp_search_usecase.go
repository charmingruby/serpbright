package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
)

func (s *ScrapperUseCaseRegistry) ProcessSerpSearchUseCase(input dto.ProcessSerpSearchInputDTO) (dto.ProcessSerpSearchOutputDTO, error) {
	searchResult, err := s.Serp.Search(input.CampaignTask)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	body, _ := json.Marshal(searchResult)

	path := fmt.Sprintf("./tmp/bright_data_%s_response.json", time.Time.Format(time.Now(), "2006-01-02 15:04:05"))
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	//	if err := s.SearchResultRepository.Store(processedResult); err != nil {
	//		return dto.ProcessSerpSearchOutputDTO{}, errors.New("Serp result insertion error: " + err.Error())
	//}

	return dto.ProcessSerpSearchOutputDTO{
		SearchResult: searchResult,
	}, nil
}
