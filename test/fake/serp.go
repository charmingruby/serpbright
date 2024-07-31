package fake

import "github.com/charmingruby/serpright/internal/scrapper/domain/entity"

func NewFakeSerp() FakeSerp {
	return FakeSerp{}
}

type FakeSerp struct{}

func (s *FakeSerp) Search(campaigntask entity.CampaignTask) entity.RawSearchData {
	return entity.RawSearchData{}
}
