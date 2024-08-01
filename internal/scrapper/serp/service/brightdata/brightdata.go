package brightdata

import (
	"fmt"

	"github.com/charmingruby/serpright/config"
)

func NewBrightData(cfg config.Config) *BrightData {
	formattedUsername := fmt.Sprintf("brd-customer-%s-zone-%s",
		cfg.BrightDataConfig.CustomerID,
		cfg.BrightDataConfig.Zone,
	)

	formattedProxyURL := fmt.Sprintf("http://%s:%s@%s:%d",
		formattedUsername,
		cfg.BrightDataConfig.Password,
		cfg.BrightDataConfig.Host,
		cfg.BrightDataConfig.Port,
	)

	return &BrightData{
		Host:       cfg.BrightDataConfig.Host,
		Port:       cfg.BrightDataConfig.Port,
		CustomerID: cfg.BrightDataConfig.CustomerID,
		Zone:       cfg.BrightDataConfig.Zone,
		Username:   formattedUsername,
		Password:   cfg.BrightDataConfig.Password,
		ProxyURL:   formattedProxyURL,
	}
}

type BrightData struct {
	Host       string
	Port       int
	CustomerID string
	Zone       string
	Username   string
	Password   string
	ProxyURL   string
}
