package factory

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

func MakeCampaignTask() entity.CampaignTask {
	return entity.CampaignTask{
		ID: "Campaign ID",
	}
}
