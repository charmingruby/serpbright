package brightdata

import (
	"fmt"

	"github.com/charmingruby/serpright/config"
)

func NewBrightData(cfg config.Config) *BrightData {
	formattedProxyURL := fmt.Sprintf("http://%s:%s@%s:%d",
		cfg.BrightDataConfig.Username,
		cfg.BrightDataConfig.Password,
		cfg.BrightDataConfig.Host,
		cfg.BrightDataConfig.Port,
	)

	return &BrightData{
		Host:         cfg.BrightDataConfig.Host,
		Port:         cfg.BrightDataConfig.Port,
		Username:     cfg.BrightDataConfig.Username,
		Password:     cfg.BrightDataConfig.Password,
		ProxyURL:     formattedProxyURL,
		DebugMode:    cfg.DebugMode,
		SearchConfig: cfg.SearchConfig,
	}
}

type BrightData struct {
	Host         string
	Port         int
	Username     string
	Password     string
	ProxyURL     string
	DebugMode    bool
	SearchConfig config.SearchConfig
}
