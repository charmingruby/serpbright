package fake

import (
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func NewFakeSerp() FakeSerp {
	return FakeSerp{}
}

type FakeSerp struct{}

func (s *FakeSerp) Search(campaigntask entity.CampaignTask) (process_entity.RawSearchData, error) {
	return process_entity.RawSearchData{}, nil
}
