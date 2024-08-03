package dto

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

type ProcessSerpSearchInputDTO struct {
	CampaignTask entity.CampaignTask
}

type ProcessSerpSearchOutputDTO struct {
	SearchResult process_entity.ResultantData
}
