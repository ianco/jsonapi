package rest

import (
	"fmt"
	"os"

	integration "github.com/ianco/jsonapi/integration"
	model "github.com/ianco/jsonapi/model"
)

var currentId int

var config model.ConfigData
var hlfSetup integration.BaseSetupImpl

// Give us some seed data
func init() {
	//RepoCreateUpdateConfig(model.ConfigData{DifficultyRating: 2, Startup: model.StartupData{AiNationCount: 1, StartupCash: 20, AiStartupCash: 10, AiAggressiveness: 1, StartupIndependentTown: 15, StartupRawSite: 5, DifficultyLevel: 1}})

	hlfSetup = integration.BaseSetupImpl{
		ConfigFile:      "./config_test.yaml",
		ChainID:         "mychannel",
		ChainCodeID:	 "abc123_cc",
		ChannelConfig:   "./mychannel.tx",
		ConnectEventHub: true,
	}

	if err := hlfSetup.Initialize(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	chain := hlfSetup.Chain
	//client := hlfSetup.Client

	isitinstalled, err := hlfSetup.IsInstalledChaincode(hlfSetup.ChainCodeID)
	if err != nil {
		fmt.Printf("IsInstalledChaincode return error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Is it installed?  %t", isitinstalled)
	if !isitinstalled {
		fmt.Printf("InstallAndInstantiateMyCC()!!!")
		if err := hlfSetup.InstallAndInstantiateCC(); err != nil {
			fmt.Printf("InstallAndInstantiateMyCC return error: %v", err)
			os.Exit(1)
		}
	}

	// Test Query Info - retrieve values before transaction
	fmt.Printf("QueryInfo()!!!")
	info, err := chain.QueryInfo()
	if err != nil {
		fmt.Printf("QueryInfo return error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("QueryInfo [%s]", info)

	config, err = RepoFindConfig()
	if err != nil {
		fmt.Printf("Configuration error, no configuration: %v", err)
		os.Exit(1)
	}
}

func RepoFindConfig() (model.ConfigData, error) {

	value, err := hlfSetup.QueryConfiguration()
	if err != nil {
		fmt.Printf("QueryConfiguration return error: %v", err)
		return config, err
	} else {
		fmt.Printf("QueryConfiguration() = %s", value)

		config, err = model.Json2Config(value)
		if err != nil {
			fmt.Printf("Not a valid Configuration: %v", err)
			return config, err
		}
	}

	return config, nil
}

func RepoCreateUpdateConfig(c model.ConfigData) (model.ConfigData, error) {

	str, _ := model.Config2Json(c)
	value, err := hlfSetup.UpdateConfiguration(str)
	if err != nil {
		fmt.Printf("UpdateConfiguration return error: %v", err)
		return config, err
	}
	fmt.Printf("UpdateConfiguration() = %s", value)
	fmt.Printf("UpdateConfiguration() atr = %s", str)

	config, err = model.Json2Config(str)
	if err != nil {
		fmt.Printf("Not a valid Configuration: %v", err)
		return config, err
	}

	return config, nil
}
