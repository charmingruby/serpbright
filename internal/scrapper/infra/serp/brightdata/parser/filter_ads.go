package parser

import (
	"fmt"
	"strings"

	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/constant"
)

func (p *BrightDataParser) FilterADs(bdr *data.BrightDataSearchResult) ([]data.Ad, []data.Ad) {
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

				inlineShoppingAD := data.TopPla{
					Link:         ad.AD.Link,
					Title:        ad.AD.Title,
					ReferralLink: ad.AD.Link,
				}

				bdr.TopPla = append(bdr.TopPla, inlineShoppingAD)
			}
		}
	}

	filteredTopADs := []data.Ad{}
	filteredBottomADs := []data.Ad{}
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
	AD            data.Ad
	blockPosition string
}
