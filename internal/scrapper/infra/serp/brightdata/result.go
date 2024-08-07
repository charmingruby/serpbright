package brightdata

import "time"

type BrightDataSearchResult struct {
	General struct {
		BasicView    bool      `json:"basic_view"`
		Language     string    `json:"language"`
		Mobile       bool      `json:"mobile"`
		PageTitle    string    `json:"page_title"`
		ResultsCnt   int       `json:"results_cnt"`
		SearchEngine string    `json:"search_engine"`
		SearchTime   float64   `json:"search_time"`
		SearchType   string    `json:"search_type"`
		Timestamp    time.Time `json:"timestamp"`
	} `json:"general"`
	Input struct {
		OriginalURL string `json:"original_url"`
		RequestID   string `json:"request_id"`
		UserAgent   string `json:"user_agent"`
	} `json:"input"`
	Knowledge struct {
		Description       string `json:"description"`
		DescriptionLink   string `json:"description_link"`
		DescriptionSource string `json:"description_source"`
		Facts             []struct {
			Key       string `json:"key"`
			Predicate string `json:"predicate"`
			Value     []struct {
				Link string `json:"link,omitempty"`
				Text string `json:"text"`
			} `json:"value"`
		} `json:"facts"`
		Images []struct {
			Image       string `json:"image"`
			ImageBase64 string `json:"image_base64"`
		} `json:"images"`
		Name     string `json:"name"`
		Subtitle string `json:"subtitle"`
		Summary  string `json:"summary"`
		Widgets  []struct {
			GlobalRank int `json:"global_rank"`
			Items      []struct {
				Link string `json:"link"`
				Name string `json:"name"`
				Rank int    `json:"rank"`
			} `json:"items"`
			Key       string `json:"key"`
			Predicate string `json:"predicate"`
			Rank      int    `json:"rank"`
			Type      string `json:"type"`
		} `json:"widgets"`
	} `json:"knowledge"`
	Navigation []struct {
		Href  string `json:"href"`
		Title string `json:"title"`
	} `json:"navigation"`
	Organic []struct {
		Description string `json:"description"`
		DisplayLink string `json:"display_link"`
		Extensions  []struct {
			Extended bool   `json:"extended"`
			Link     string `json:"link"`
			Rank     int    `json:"rank"`
			Text     string `json:"text"`
			Type     string `json:"type"`
		} `json:"extensions,omitempty"`
		GlobalRank int    `json:"global_rank"`
		Image      string `json:"image,omitempty"`
		ImageAlt   string `json:"image_alt,omitempty"`
		ImageURL   string `json:"image_url,omitempty"`
		Link       string `json:"link"`
		Rank       int    `json:"rank"`
		Title      string `json:"title"`
	} `json:"organic"`
	Overview struct {
		Thematic struct {
			Rating string `json:"rating"`
			Title  string `json:"title"`
			Year   string `json:"year"`
		} `json:"thematic"`
	} `json:"overview"`
	Pagination struct {
		CurrentPage   int    `json:"current_page"`
		NextPage      int    `json:"next_page"`
		NextPageLink  string `json:"next_page_link"`
		NextPageStart int    `json:"next_page_start"`
		Pages         []struct {
			Link  string `json:"link"`
			Page  int    `json:"page"`
			Start int    `json:"start"`
		} `json:"pages"`
	} `json:"pagination"`
	PeopleAlsoAsk []struct {
		AnswerDisplayLink string `json:"answer_display_link"`
		AnswerHtml        string `json:"answer_html"`
		AnswerLink        string `json:"answer_link"`
		AnswerSource      string `json:"answer_source"`
		Answers           []struct {
			Items [][]struct {
				Text string `json:"text"`
			} `json:"items,omitempty"`
			MoreLink string `json:"more_link,omitempty"`
			MoreText string `json:"more_text,omitempty"`
			Rank     int    `json:"rank"`
			Title    string `json:"title,omitempty"`
			Type     string `json:"type"`
			Value    *struct {
				Text string `json:"text"`
			} `json:"value,omitempty"`
		} `json:"answers"`
		GlobalRank   int    `json:"global_rank"`
		Question     string `json:"question"`
		QuestionLink string `json:"question_link"`
		Rank         int    `json:"rank"`
	} `json:"people_also_ask"`
	Related []struct {
		GlobalRank int    `json:"global_rank"`
		Link       string `json:"link"`
		ListGroup  bool   `json:"list_group"`
		Rank       int    `json:"rank"`
		Text       string `json:"text"`
	} `json:"related"`
	SnackPack []struct {
		Address           string   `json:"address"`
		Cid               string   `json:"cid"`
		GlobalRank        int      `json:"global_rank"`
		MapsLink          string   `json:"maps_link"`
		Name              string   `json:"name"`
		Phone             string   `json:"phone,omitempty"`
		Rank              int      `json:"rank"`
		Site              string   `json:"site,omitempty"`
		Tags              []string `json:"tags,omitempty"`
		WorkStatus        string   `json:"work_status"`
		WorkStatusDetails string   `json:"work_status_details,omitempty"`
	} `json:"snack_pack"`
	SnackPackMap struct {
		Altitude  int     `json:"altitude"`
		Image     string  `json:"image"`
		ImageAlt  string  `json:"image_alt"`
		ImageURL  string  `json:"image_url"`
		Latitude  float64 `json:"latitude"`
		Link      string  `json:"link"`
		Longitude float64 `json:"longitude"`
	} `json:"snack_pack_map"`
	TopAds []struct {
		Description string `json:"description"`
		DisplayLink string `json:"display_link"`
		Extensions  []struct {
			Description string `json:"description"`
			Extended    bool   `json:"extended"`
			Link        string `json:"link"`
			Text        string `json:"text"`
			Type        string `json:"type"`
		} `json:"extensions"`
		GlobalRank   int    `json:"global_rank"`
		Link         string `json:"link"`
		Rank         int    `json:"rank"`
		ReferralLink string `json:"referral_link"`
		Title        string `json:"title"`
	} `json:"top_ads"`
}
