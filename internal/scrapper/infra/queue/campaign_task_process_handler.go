package queue

import (
	"encoding/json"
	"log"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func NewCampaignTaskProcessHandler(scrapperService usecase.ScrapperUseCase) CampaignTaskProcessHandler {
	return CampaignTaskProcessHandler{
		ScrapperService: scrapperService,
	}
}

type CampaignTaskProcessHandler struct {
	ScrapperService usecase.ScrapperUseCase
}

func (h *CampaignTaskProcessHandler) Handle(msg []byte) {
	var ucInput dto.ProcessSerpSearchInputDTO
	if err := json.Unmarshal(msg, &ucInput); err != nil {
		log.Fatal(err)
	}

	output, err := h.ScrapperService.ProcessSerpSearchUseCase(ucInput)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PROCESSED DATA: %v", output)
}
