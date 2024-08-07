package queue

import (
	"encoding/json"
	"log"

	"github.com/charmingruby/serpright/internal/scrapper/domain/dto"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
)

func NewCampaignTaskProcessHandler(scrapperService usecase.ScrapperUseCase, debugMode bool) CampaignTaskProcessHandler {
	return CampaignTaskProcessHandler{
		ScrapperService: scrapperService,
		DebugMode:       debugMode,
	}
}

type CampaignTaskProcessHandler struct {
	ScrapperService usecase.ScrapperUseCase
	DebugMode       bool
}

func (h *CampaignTaskProcessHandler) Handle(msg []byte) {
	var task entity.CampaignTask
	if err := json.Unmarshal(msg, &task); err != nil {
		log.Fatal(err)
	}

	input := dto.ProcessSerpSearchInputDTO{
		CampaignTask: task,
	}

	output, err := h.ScrapperService.ProcessSerpSearchUseCase(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PROCESSED DATA: %v", output)
}
