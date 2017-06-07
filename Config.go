package main


type StartupData struct {
    AiNationCount          uint8  `json:"ai_nation_count"`
    StartupCash            uint8  `json:"start_up_cash"`
    AiStartupCash          uint8  `json:"ai_start_up_cash"`
    AiAggressiveness       uint8  `json:"ai_aggressiveness"`
    StartupIndependentTown uint8  `json:"start_up_independent_town"`
    StartupRawSite         uint8  `json:"start_up_raw_site"`
    DifficultyLevel        uint8  `json:"difficulty_level"`
}

type ConfigData struct {
    DifficultyRating uint8       `json:"difficulty_rating"`
    Startup          StartupData `json:"startup"`
}



