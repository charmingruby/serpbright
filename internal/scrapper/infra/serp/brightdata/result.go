package brightdata

import "time"

type BrightDataSearchResult struct {
	General    GeneralInfo  `json:"general"`
	Images     []Image      `json:"images"`
	Input      InputInfo    `json:"input"`
	Navigation []Navigation `json:"navigation"`
	Organic    []Organic    `json:"organic"`
	Overview   Overview     `json:"overview"`
	Pagination Pagination   `json:"pagination"`
	Related    []Related    `json:"related"`
	TopAds     []Ad         `json:"top_ads"`
	TopPla     []TopPla     `json:"top_pla"`
	BottomAds  []Ad         `json:"bottom_ads"`
	HTML       string       `json:"html"`
}

type Ad struct {
	Description  string      `json:"description"`
	DisplayLink  string      `json:"display_link"`
	Extensions   []Extension `json:"extensions,omitempty"`
	GlobalRank   int         `json:"global_rank"`
	Link         string      `json:"link"`
	Rank         int         `json:"rank"`
	ReferralLink string      `json:"referral_link"`
	Title        string      `json:"title"`
}

type TopPla struct {
	GlobalRank   int     `json:"global_rank"`
	Image        string  `json:"image,omitempty"`
	ImageAlt     string  `json:"image_alt,omitempty"`
	ImageBase64  string  `json:"image_base64,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	Link         string  `json:"link,omitempty"`
	Price        string  `json:"price,omitempty"`
	Rating       float64 `json:"rating,omitempty"`
	ReviewsCnt   int     `json:"reviews_cnt,omitempty"`
	Rank         int     `json:"rank"`
	ReferralLink string  `json:"referral_link,omitempty"`
	Shop         string  `json:"shop,omitempty"`
	Title        string  `json:"title"`
	ViewAll      bool    `json:"view_all,omitempty"`
}

type Extension struct {
	Link string `json:"link"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type GeneralInfo struct {
	BasicView    bool      `json:"basic_view"`
	Language     string    `json:"language"`
	Mobile       bool      `json:"mobile"`
	PageTitle    string    `json:"page_title"`
	ResultsCnt   int       `json:"results_cnt"`
	SearchEngine string    `json:"search_engine"`
	SearchTime   float64   `json:"search_time"`
	SearchType   string    `json:"search_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type Image struct {
	GlobalRank    int    `json:"global_rank"`
	Image         string `json:"image"`
	ImageAlt      string `json:"image_alt"`
	ImageURL      string `json:"image_url"`
	Link          string `json:"link"`
	OriginalImage string `json:"original_image"`
	Rank          int    `json:"rank"`
	Source        string `json:"source"`
	SourceLogo    string `json:"source_logo"`
	Title         string `json:"title"`
}

type InputInfo struct {
	OriginalURL string `json:"original_url"`
	RequestID   string `json:"request_id"`
}

type Navigation struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type Organic struct {
	Description string      `json:"description"`
	DisplayLink string      `json:"display_link"`
	Duration    string      `json:"duration,omitempty"`
	DurationSec int         `json:"duration_sec,omitempty"`
	Extensions  []Extension `json:"extensions,omitempty"`
	GlobalRank  int         `json:"global_rank"`
	Image       string      `json:"image,omitempty"`
	ImageAlt    string      `json:"image_alt,omitempty"`
	ImageURL    string      `json:"image_url,omitempty"`
	Link        string      `json:"link"`
	Rank        int         `json:"rank"`
	Title       string      `json:"title"`
}

type Overview struct {
	Kgmid string `json:"kgmid"`
	Title string `json:"title"`
}

type Pagination struct {
	CurrentPage   int    `json:"current_page"`
	NextPage      int    `json:"next_page"`
	NextPageLink  string `json:"next_page_link"`
	NextPageStart int    `json:"next_page_start"`
	Pages         []Page `json:"pages"`
}

type Page struct {
	Link  string `json:"link"`
	Page  int    `json:"page"`
	Start int    `json:"start"`
}

type Related struct {
	GlobalRank int    `json:"global_rank"`
	Link       string `json:"link"`
	ListGroup  bool   `json:"list_group"`
	Rank       int    `json:"rank"`
	Text       string `json:"text"`
}
