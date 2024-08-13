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
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata/data"
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
}

func (s *BrightData) doRequest(reqURL string) (data.BrightDataSearchResult, error) {
	proxy, err := url.Parse(s.ProxyURL)
	if err != nil {
		slog.Error("Proxy URL parse error: ")
		return data.BrightDataSearchResult{}, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", reqURL+"&brd_json=html", nil)
	if err != nil {
		slog.Error("Request creation error: ")
		return data.BrightDataSearchResult{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request error: ")
		return data.BrightDataSearchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Request failed with status: %d", resp.StatusCode))
		return data.BrightDataSearchResult{}, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: ")
		return data.BrightDataSearchResult{}, err
	}

	if s.DebugMode {
		path := fmt.Sprintf("./tmp/bright_data_%s_response.json", time.Time.Format(time.Now(), "2006-01-02 15:04:05"))
		err := os.WriteFile(path, body, 0644)
		if err != nil {
			return data.BrightDataSearchResult{}, err
		}
	}

	var serpResult data.BrightDataSearchResult
	if err := json.Unmarshal(body, &serpResult); err != nil {
		slog.Error("Decode error: ")
		return data.BrightDataSearchResult{}, err
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
		Device:       s.extractDeviceFromTask(campaignTask),
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
	}

	url := baseURL + strings.Join(builtParams[:], "&") + defaultParams

	return url
}
func (s *BrightData) extractDeviceFromTask(task entity.CampaignTask) string {
	if task.Device == constant.MobileDevice {
		if task.MobileType == constant.MobileTypeAndroid {
			return AndroidDevice
		}

		return IOSDevice
	}

	return DesktopDevice
}
