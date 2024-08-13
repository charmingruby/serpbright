package brightdata

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/constant"
)

func BrighDataResultToSearchResult(result BrightDataSearchResult, task entity.CampaignTask) process_entity.SearchResult {
	return process_entity.SearchResult{
		ID:              "",
		Results:         []process_entity.SearchResultItem{},
		ShoppingResults: []process_entity.ShoppingSearchResultItem{},
		Task:            task,
		SearchUrl:       result.General.SearchType,
		HTMLData:        result.HTML,
		CreatedAt:       time.Now(),
	}
}

func (s *BrightData) filterADs(bdr *BrightDataSearchResult) ([]Ad, []Ad) {
	var adsToFilter []ADFilter

	if len(bdr.TopAds) != 0 {
		for _, ad := range bdr.TopAds {
			adsToFilter = append(adsToFilter, ADFilter{
				AD:            ad,
				blockPosition: "top",
			})
		}
	}

	if len(bdr.BottomAds) != 0 {
		for _, ad := range bdr.BottomAds {
			adsToFilter = append(adsToFilter, ADFilter{
				AD:            ad,
				blockPosition: "bottom",
			})
		}
	}

	var adDuplicatedOrEmpty bool
	for _, ad := range adsToFilter {
		adDuplicatedOrEmpty = false

		if ad.AD.Link == "" {
			continue
		}

		if strings.HasPrefix(ad.AD.ReferralLink, "/aclk") {
			for _, shoppingItem := range bdr.TopPla {
				if ad.AD.Link == shoppingItem.Link {
					adDuplicatedOrEmpty = true
					break
				}
			}

			if !adDuplicatedOrEmpty {
				ad.AD.ReferralLink = fmt.Sprintf("https://www.%s%s",
					constant.GoogleDomain,
					ad.AD.ReferralLink,
				)

				inlineShoppingAD := TopPla{
					Link:         ad.AD.Link,
					Title:        ad.AD.Title,
					ReferralLink: ad.AD.Link,
				}

				bdr.TopPla = append(bdr.TopPla, inlineShoppingAD)
			}
		}
	}

	filteredTopADs := []Ad{}
	filteredBottomADs := []Ad{}
	for _, ad := range adsToFilter {
		if ad.blockPosition == "top" {
			filteredTopADs = append(filteredTopADs, ad.AD)
			continue
		}

		filteredBottomADs = append(filteredBottomADs, ad.AD)
	}

	return filteredTopADs, filteredBottomADs
}

type ADFilter struct {
	AD            Ad
	blockPosition string
}
