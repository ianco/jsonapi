package rest

import (
	model "github.com/ianco/jsonapi/model"
)

var currentId int

var config model.ConfigData

// Give us some seed data
func init() {
	RepoCreateUpdateConfig(model.ConfigData{DifficultyRating: 2, Startup: model.StartupData{AiNationCount: 1, StartupCash: 20, AiStartupCash: 10, AiAggressiveness: 1, StartupIndependentTown: 15, StartupRawSite: 5, DifficultyLevel: 1}})
}

func RepoFindConfig() model.ConfigData {
	return config
}

func RepoCreateUpdateConfig(c model.ConfigData) model.ConfigData {
	config = c
	return config
}
