package brightdata

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func (s *BrightData) ExecSearch(svc usecase.ScrapperUseCase) (dto.ProcessSerpSearchOutputDTO, error) {
	op, err := svc.ProcessSerpSearchUseCase(dto.ProcessSerpSearchInputDTO{CampaignTaskID: "id"})
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	return op, nil
}
