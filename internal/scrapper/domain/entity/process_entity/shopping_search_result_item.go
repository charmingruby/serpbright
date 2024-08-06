package process_entity

type ShoppingSearchResultItem struct {
	SearchResultItem   `bson:",inline"`
	CompanyName        string  `json:"companyName" bson:"companyName,omitempty"`
	RegularPriceString string  `json:"regularPriceString" bson:"regularPriceString,omitempty"`
	PriceString        string  `json:"priceString" bson:"priceString,omitempty"`
	Price              float64 `json:"price" bson:"price,omitempty"`
	RegularPrice       float64 `json:"regularPrice" bson:"regularPrice,omitempty"`
	Currency           string  `json:"currency" bson:"currency,omitempty"`
	Channel            string  `json:"channel" bson:"-"`
}
