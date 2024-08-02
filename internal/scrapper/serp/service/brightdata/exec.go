package brightdata

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
	"github.com/charmingruby/serpright/test/factory"
)

func (s *BrightData) ExecSearch(svc usecase.ScrapperUseCase) (dto.ProcessSerpSearchOutputDTO, error) {
	input := dto.ProcessSerpSearchInputDTO{
		CampaignTask: factory.MakeCampaignTask(),
	}

	op, err := svc.ProcessSerpSearchUseCase(input)
	if err != nil {
		return dto.ProcessSerpSearchOutputDTO{}, err
	}

	return op, nil
}
