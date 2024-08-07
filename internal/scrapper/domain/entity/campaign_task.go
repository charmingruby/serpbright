package entity

import "time"

type CampaignTask struct {
	ID                  string      `json:"id" bson:"_id"` // CampaignId;Device;MobileType;Page;Keyword;GeoLocation;CreatedAtUTC
	CreatedAt           TimeWrapper `json:"createdAt" bson:"createdAt"`
	SearchType          string      `json:"searchType" bson:"searchType"`
	BrokerId            string      `json:"brokerId" bson:"brokerId"`
	GeoLocation         string      `json:"geoLocation" bson:"geoLocation"`
	SearchEngineDomain  string      `json:"searchEngineDomain,omitempty" bson:"searchEngineDomain,omitempty"`
	LocaleCountry       string      `json:"localeCountry,omitempty" bson:"localeCountry,omitempty"`
	Locale              string      `json:"locale,omitempty" bson:"locale,omitempty"`
	Device              string      `json:"device,omitempty" bson:"device,omitempty"` // Desktop, mobile, tablet
	MobileType          string      `json:"mobileType" bson:"mobileType"`             // ipad (device tablet), iphone (device mobile), android (device tablet ou mobile)
	Page                uint8       `json:"page" bson:"page"`
	BrandName           string      `json:"brandName" bson:"brandName"`
	Keyword             string      `json:"keyword" bson:"keyword"`
	Domain              string      `json:"domain" bson:"domain"`
	ShouldAlert         bool        `json:"shouldAlert" bson:"shouldAlert"`
	AlertType           string      `json:"alertType" bson:"alertType"`
	SendReport          bool        `json:"sendReport" bson:"sendReport"`
	ReportType          string      `json:"reportType" bson:"reportType"`
	ReportTemplate      string      `json:"reportTemplate" bson:"reportTemplate"`
	MaxQueueWaitMinutes uint8       `json:"maxQueueWaitMinutes" bson:"maxQueueWaitMinutes"`
	StartedAt           TimeWrapper `json:"startedAt" bson:"startedAt"`
	FinishedAt          TimeWrapper `json:"finishedAt" bson:"finishedAt"`
	SearchUrl           string      `json:"searchUrl" bson:"searchUrl"`
	HtmlDataUrl         string      `json:"htmlDataUrl" bson:"htmlDataUrl"`
}

type TimeWrapper struct {
	Date time.Time `json:"$date"`
}
