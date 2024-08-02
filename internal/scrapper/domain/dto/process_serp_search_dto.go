package dto

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

type ProcessSerpSearchInputDTO struct {
	CampaignTask entity.CampaignTask
}

type ProcessSerpSearchOutputDTO struct {
	SearchResult entity.ResultantData
}
