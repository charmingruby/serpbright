package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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
		log.Fatal(err.Error())
	}

	body, err := json.Marshal(output.SearchResult)
	if err != nil {
		slog.Error(err.Error())
	}

	path := fmt.Sprintf("./tmp/results/search_result_%s_response.json", time.Time.Format(time.Now(), "2006-01-02 15:04:05"))
	if err := os.WriteFile(path, body, 0644); err != nil {
		slog.Error(err.Error())
	}

	log.Printf("TASK PROCESSED: %v", output.SearchResult.Task.BrandName)
}
