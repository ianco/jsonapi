package main


type StartupData struct {
    AiNationCount          string `json:"ai_nation_count"`
    StartupCash            string `json:"start_up_cash"`
    AiStartupCash          string `json:"ai_start_up_cash"`
    AiAggressiveness       string `json:"ai_aggressiveness"`
    StartupIndependentTown string `json:"start_up_independent_town"`
    StartupRawSite         string `json:"start_up_raw_site"`
    DifficultyLevel        string `json:"difficulty_level"`
}

type ConfigData struct {
    DifficultyRating string      `json:"difficulty_rating"`
    Startup          StartupData `json:"startup"`
}



