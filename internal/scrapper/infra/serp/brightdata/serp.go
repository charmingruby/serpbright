package brightdata

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/charmingruby/serpright/internal/common/helper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *BrightData) Search(campaigntask entity.CampaignTask) (process_entity.SearchResult, error) {
	proxy, err := url.Parse(s.ProxyURL)
	if err != nil {
		slog.Error("Proxy URL parse error: " + err.Error())
		return process_entity.SearchResult{}, err
	}

	slog.Info(fmt.Sprintf("Using proxy URL: %s", s.ProxyURL))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", "https://www.google.com/search?q=pizza&brd_json=1", nil)
	if err != nil {
		slog.Error("Request creation error: " + err.Error())
		return process_entity.SearchResult{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request error: " + err.Error())
		return process_entity.SearchResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error(fmt.Sprintf("Request failed with status: %d. Response body: %s", resp.StatusCode, string(body)))
		return process_entity.SearchResult{}, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: " + err.Error())
		return process_entity.SearchResult{}, err
	}

	var serpResult BrightDataSearchResult
	if err := json.Unmarshal(body, &serpResult); err != nil {
		slog.Error("Decode error: " + err.Error())
		return process_entity.SearchResult{}, err
	}

	if debugMode {
		if err := helper.DebugJSON(serpResult); err != nil {
			return process_entity.SearchResult{}, err
		}
	}

	rawData := BrighDataResultToSearchResult(serpResult)

	return rawData, nil
}
