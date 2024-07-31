package dto

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

type ProcessSerpSearchInputDTO struct {
	CampaignTaskID string
}

type ProcessSerpSearchOutputDTO struct {
	SearchResult entity.ResultantData
}
