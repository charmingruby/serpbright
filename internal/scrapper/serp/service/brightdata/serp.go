package brightdata

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/charmingruby/serpright/internal/scrapper/domain/entity"
	"github.com/charmingruby/serpright/internal/scrapper/domain/entity/process_entity"
)

func (s *BrightData) Search(campaigntask entity.CampaignTask) (process_entity.RawSearchData, error) {
	proxy, err := url.Parse(s.ProxyURL)
	if err != nil {
		slog.Error("Proxy URL parse error: " + err.Error())
		return process_entity.RawSearchData{}, err
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

	url := "https://www.google.com/search?q=pizza"

	resp, err := client.Get(url)
	if err != nil {
		slog.Error("Request error: " + err.Error())
		return process_entity.RawSearchData{}, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Error("Decode error: " + err.Error())
		return process_entity.RawSearchData{}, err
	}

	json, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return process_entity.RawSearchData{}, err
	}

	fmt.Println(string(json))

	return process_entity.RawSearchData{}, nil
}
