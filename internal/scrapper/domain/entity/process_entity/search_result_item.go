package process_entity

import "time"

// Ocorrência em resultado de buscas, "único" por busca
type SearchResultItem struct {
	Id                 *string   `json:"id" bson:"_id,omitempty"`
	CampaignTaskId     string    `json:"campaignTaskId" bson:"campaignTaskId,omitempty"`
	CampaignId         string    `json:"campaignId" bson:"campaignId,omitempty"`
	Device             string    `json:"device" bson:"device,omitempty"`
	DeviceType         string    `json:"deviceType" bson:"deviceType,omitempty"`
	Keyword            string    `json:"keyword" bson:"keyword,omitempty"`
	Page               string    `json:"page" bson:"page,omitempty"`
	GeoLocation        string    `json:"geoLocation" bson:"geoLocation,omitempty"`
	Position           uint8     `json:"position" bson:"position,omitempty"`
	PositionOnPage     string    `json:"positionOnPage" bson:"positionOnPage,omitempty"`
	Title              string    `json:"title" bson:"title,omitempty"`
	BrandInTitle       bool      `json:"brandInTitle" bson:"brandInTitle,omitempty"`
	BrandInUrl         bool      `json:"brandInUrl" bson:"brandInUrl,omitempty"`
	BrandInDescription bool      `json:"brandInDescription" bson:"brandInDescription,omitempty"`
	Url                string    `json:"url" bson:"url,omitempty"`
	UrlSequence        []string  `json:"urlSequence" bson:"urlSequence,omitempty"`
	SiteType           string    `json:"siteType" bson:"siteType,omitempty"`
	TrackingUrl        string    `json:"trackingUrl" bson:"trackingUrl,omitempty"`
	DisplayedDomain    string    `json:"displayedDomain" bson:"displayedDomain,omitempty"`
	Domain             string    `json:"domain" bson:"domain,omitempty"`
	Phone              string    `json:"phone" bson:"phone,omitempty"`
	IsOwner            bool      `json:"isOwner" bson:"isOwner,omitempty"`
	SubDomain          string    `json:"subDomain" bson:"subDomain,omitempty"`
	LinkUrl            string    `json:"linkUrl" bson:"linkUrl,omitempty"`
	DisplayedUrl       string    `json:"displayedUrl" bson:"displayedUrl,omitempty"`
	Description        string    `json:"description" bson:"description,omitempty"`
	Type               string    `json:"type" bson:"type,omitempty"` // ad, organic, video, post, etc
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	RedirectSeconds    float64   `json:"redirectSeconds" bson:"redirectSeconds,omitempty"`
	RedirectHTTPCode   int       `json:"redirectHttpCode" bson:"redirectHttpCode,omitempty"`
	Channel            string    `json:"channel" bson:"channel,omitempty"`
	Evidence           string    `json:"evidence" bson:"evidence,omitempty"`
	HasError           bool
	//Sitelinks     Sitelinks `json:"sitelinks,omitempty"`
}
