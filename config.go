package main

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	List ListConfig `json:"list"`
	UI   UIConfig   `json:"ui"`
}

type ListConfig struct {
	Running  bool   `json:"running"`
	Name     string `json:"name"`
	JSON     bool   `json:"json"`
	LogsTail string `json:"logsTail"`
}

type UIConfig struct {
	FooterColor string `json:"footerColor"`
	StatusColor string `json:"statusColor"`
	AccentColor string `json:"accentColor"`
}

func defaultConfig() AppConfig {
	return AppConfig{
		List: ListConfig{
			Running:  false,
			Name:     "",
			JSON:     false,
			LogsTail: "100",
		},
		UI: UIConfig{
			FooterColor: "white",
			StatusColor: "white",
			AccentColor: "yellow",
		},
	}
}

func loadConfig(path string) AppConfig {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, &cfg)

	if cfg.List.LogsTail == "" {
		cfg.List.LogsTail = "100"
	}

	if cfg.UI.FooterColor == "" {
		cfg.UI.FooterColor = "white"
	}
	if cfg.UI.StatusColor == "" {
		cfg.UI.StatusColor = "white"
	}
	if cfg.UI.AccentColor == "" {
		cfg.UI.AccentColor = "yellow"
	}

	return cfg
}