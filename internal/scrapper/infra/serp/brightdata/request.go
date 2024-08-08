package brightdata

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmingruby/serpright/internal/common/helper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/constant"
)

const (
	baseURL       = "https://www.google.com/search?"
	defaultParams = "&output=json&exclude_fields=inline_videos,inline_tweets,related_questions,knowledge_graph,search_parameters,pagination,organic_results"
)

type brightDataRequestParams struct {
	UULE         string
	GoogleDomain string
	GL           string
	HL           string
	Q            string
	Device       string
	Page         int
	IncludeHTML  bool
}

func (s *BrightData) doHTMLRequest(reqURL string) (string, error) {
	proxy, err := url.Parse(s.ProxyURL)
	if err != nil {
		slog.Error("Proxy URL parse error: " + err.Error())
		return "", err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		slog.Error("Request creation error: " + err.Error())
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request error: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error(fmt.Sprintf("Request failed with status: %d. Response body: %s", resp.StatusCode, string(body)))
		return "", fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: " + err.Error())
		return "", err
	}

	if s.DebugMode {
		path := fmt.Sprintf("./tmp/bright_data_%s_response.html", time.Time.Format(time.Now(), "2006-01-02 15:04:05"))
		err := os.WriteFile(path, body, 0644)
		if err != nil {
			return "", err
		}
	}

	return string(body), nil
}

func (s *BrightData) doJSONRequest(reqURL string) (BrightDataSearchResult, error) {
	proxy, err := url.Parse(s.ProxyURL)
	if err != nil {
		slog.Error("Proxy URL parse error: " + err.Error())
		return BrightDataSearchResult{}, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", reqURL+"&brd_json=1", nil)
	if err != nil {
		slog.Error("Request creation error: " + err.Error())
		return BrightDataSearchResult{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request error: " + err.Error())
		return BrightDataSearchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error(fmt.Sprintf("Request failed with status: %d. Response body: %s", resp.StatusCode, string(body)))
		return BrightDataSearchResult{}, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: " + err.Error())
		return BrightDataSearchResult{}, err
	}

	if s.DebugMode {
		path := fmt.Sprintf("./tmp/bright_data_%s_response.json", time.Time.Format(time.Now(), "2006-01-02 15:04:05"))
		err := os.WriteFile(path, body, 0644)
		if err != nil {
			return BrightDataSearchResult{}, err
		}
	}

	var serpResult BrightDataSearchResult
	if err := json.Unmarshal(body, &serpResult); err != nil {
		slog.Error("Decode error: " + err.Error())
		return BrightDataSearchResult{}, err
	}

	return serpResult, nil
}

func (s *BrightData) buildBrightDataRequestURL(campaignTask entity.CampaignTask) string {
	base64GeoLocation := base64.StdEncoding.EncodeToString([]byte(campaignTask.GeoLocation))
	itemsPerPage := 10

	params := brightDataRequestParams{
		UULE:         url.QueryEscape("w+CAIQICI" + constant.UULEKeys[len(campaignTask.GeoLocation)] + base64GeoLocation),
		GoogleDomain: helper.EmptyString(campaignTask.SearchEngineDomain, constant.GoogleDomain),
		GL:           helper.EmptyString(campaignTask.LocaleCountry, "br"),
		HL:           helper.EmptyString(campaignTask.Locale, "pt-br"),
		Q:            url.QueryEscape(campaignTask.Keyword),
		IncludeHTML:  s.IncludeHTML,
		Device:       "0",
		Page:         int(campaignTask.Page) * itemsPerPage,
	}

	builtParams := []string{
		"google_domain=" + params.GoogleDomain,
		"uule=" + params.UULE,
		"gl=" + params.GL,
		"hl=" + params.HL,
		"q=" + params.Q,
		"brd_mobile=" + params.Device,
		"page=" + strconv.Itoa(int(params.Page)),
		"include_html=" + strconv.FormatBool(params.IncludeHTML),
	}

	url := baseURL + strings.Join(builtParams[:], "&") + defaultParams

	return url
}
