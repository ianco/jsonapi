package model

import (
	"encoding/json"
)

type StartupData struct {
	AiNationCount          uint8 `json:"ai_nation_count"`
	StartupCash            uint8 `json:"start_up_cash"`
	AiStartupCash          uint8 `json:"ai_start_up_cash"`
	AiAggressiveness       uint8 `json:"ai_aggressiveness"`
	StartupIndependentTown uint8 `json:"start_up_independent_town"`
	StartupRawSite         uint8 `json:"start_up_raw_site"`
	DifficultyLevel        uint8 `json:"difficulty_level"`
}

type ConfigData struct {
	DifficultyRating uint8       `json:"difficulty_rating"`
	Startup          StartupData `json:"startup"`
}

// Config2Json converts a ConfigData object to Json encoding
func Config2Json(c ConfigData) (string, error) {
	var s string
	var byt []byte
	byt, err := json.Marshal(c)
	if err == nil {
		s = string(byt)
	}
	return s, err
}

// Json2Config converts a Json object and returns a ConfigData object
func Json2Config(s string) (ConfigData, error) {
	var c ConfigData
	var byt []byte
	byt = []byte(s)
	err := json.Unmarshal(byt, &c)
	return c, err
}

