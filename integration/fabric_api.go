/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package integration

import (
	"fmt"
	"time"
	"os"
	"path"

	"github.com/hyperledger/fabric-sdk-go/config"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"

	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fcutil "github.com/hyperledger/fabric-sdk-go/fabric-client/util"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	logging "github.com/op/go-logging"
)

// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	Client             fabricClient.Client
	OrdererAdminClient fabricClient.Client
	Chain              fabricClient.Chain
	EventHub           events.EventHub
	ConnectEventHub    bool
	ConfigFile         string
	ChainID            string
	ChainCodeID        string
	Initialized        bool
	ChannelConfig      string
}

// Logger this is a comment to fool linter
var Logger = logging.MustGetLogger("BaseSetup")

// Initialize reads configuration from file and sets up client, chain and event hub
func (setup *BaseSetupImpl) Initialize() error {
	Logger.Infof("Call InitConfig()")
	if err := setup.InitConfig(); err != nil {
		return fmt.Errorf("Init from config failed: %v", err)
	}

	Logger.Infof("Call InitFactories()")
	// Initialize bccsp factories before calling get client
	err := bccspFactory.InitFactories(config.GetCSPConfig())
	if err != nil {
		return fmt.Errorf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}

	Logger.Infof("Call GetClient()")
	client, err := fcutil.GetClient("admin", "adminpw", "/tmp/enroll_user")
	if err != nil {
		return fmt.Errorf("Create client failed: %v", err)
	}
	//clientUser := client.GetUserContext()
	Logger.Infof("Client: %+v", client)
	Logger.Infof("Client.cryptoSuite: %+v", client.GetCryptoSuite())
	Logger.Infof("Client.stateStore: %+v", client.GetStateStore())
	//Logger.Infof("Client.userContext: %+v", client.GetUserContext())

	setup.Client = client

	Logger.Infof("Call GetAdmin()")
	org1Admin, err := GetAdmin(client, "org1")
	if err != nil {
		return fmt.Errorf("Error getting org admin user: %v", err)
	}

	chain, err := fcutil.GetChain(setup.Client, setup.ChainID)
	if err != nil {
		return fmt.Errorf("Create chain (%s) failed: %v", setup.ChainID, err)
	}
	setup.Chain = chain

	ordererAdmin, err := GetOrdererAdmin(client)
	if err != nil {
		return fmt.Errorf("Error getting orderer admin user: %v", err)
	}

	Logger.Infof("******* CreateAndJoinChannel()")
	// Create and join channel
	if err := fcutil.CreateAndJoinChannel(client, ordererAdmin, org1Admin, chain, setup.ChannelConfig); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	Logger.Infof("******* SetUserContext()")
	client.SetUserContext(org1Admin)
	Logger.Infof("******* setupEventHub()")
	if err := setup.setupEventHub(client); err != nil {
		return err
	}

	setup.Initialized = true

	return nil
}

func (setup *BaseSetupImpl) setupEventHub(client fabricClient.Client) error {
	Logger.Infof("******* getEventHub()")
	eventHub, err := getEventHub(client)
	if err != nil {
		return err
	}

	if setup.ConnectEventHub {
		Logger.Infof("******* Connect()")
		if err := eventHub.Connect(); err != nil {
			return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
		}
		Logger.Infof("connected to eventHub")
	}
	setup.EventHub = eventHub

	return nil
}

// InitConfig ...
func (setup *BaseSetupImpl) InitConfig() error {
	Logger.Infof("Call InitConfig(%s)")
	if err := config.InitConfig(setup.ConfigFile); err != nil {
		return err
	}
	return nil
}

// InstantiateCC ...
func (setup *BaseSetupImpl) InstantiateCC(chainCodeID string, chainID string, chainCodePath string, chainCodeVersion string, args []string) error {
	if err := fcutil.SendInstantiateCC(setup.Chain, chainCodeID, chainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, setup.EventHub); err != nil {
		return err
	}
	return nil
}

// InstallCC ...
func (setup *BaseSetupImpl) InstallCC(chainCodeID string, chainCodePath string, chainCodeVersion string, chaincodePackage []byte) error {
	if err := fcutil.SendInstallCC(setup.Client, setup.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, setup.Chain.GetPeers(), setup.GetDeployPath()); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}
	return nil
}

// GetDeployPath ..
func (setup *BaseSetupImpl) GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../chaincode_1")
}

// InstallAndInstantiateCC ..
func (setup *BaseSetupImpl) InstallAndInstantiateCC() error {

	chainCodePath := "cc"
	chainCodeVersion := "v0"

	if setup.ChainCodeID == "" {
		setup.ChainCodeID = fcutil.GenerateRandomID()
	}

	Logger.Infof("InstallCC()!!!")
	if err := setup.InstallCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, nil); err != nil {
		return err
	}

	var args []string
	var s1 string
	s1 = "{\"difficulty_rating\": 4, \"startup\": {\"ai_nation_count\": 3, \"start_up_cash\": 10, \"ai_start_up_cash\": 20, \"ai_aggressiveness\": 3, \"start_up_independent_town\": 15, \"start_up_raw_site\": 55, \"difficulty_level\": 4 }}"
	args = append(args, "init")
	args = append(args, s1)
	args = append(args, "alice")
	args = append(args, "bob")

	Logger.Infof("InstantiateCC()!!!")
	return setup.InstantiateCC(setup.ChainCodeID, setup.ChainID, chainCodePath, chainCodeVersion, args)
}

// Query ...
func (setup *BaseSetupImpl) Query(chainID string, chainCodeID string, args []string) (string, error) {
	Logger.Infof("In Query()")
	Logger.Infof("Chain       :%+v", setup.Chain)
	Logger.Infof("ChainCodeID :%+v", chainCodeID)
	Logger.Infof("ChainID     :%+v", chainID)
	Logger.Infof("Primary peer:%+v", setup.Chain.GetPrimaryPeer())
	transactionProposalResponses, _, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, chainCodeID, chainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, nil)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	return string(transactionProposalResponses[0].GetResponsePayload()), nil
}

// QueryMyAsset ...
func (setup *BaseSetupImpl) QueryConfiguration() (string, error) {

	var args []string
	//args = append(args, "invoke")
	args = append(args, "query_config")
	return setup.Query(setup.ChainID, setup.ChainCodeID, args)
}

// UpdateMyAsset ...
func (setup *BaseSetupImpl) UpdateConfiguration(s2 string) (string, error) {

	var args []string
	//s2 = "{\"difficulty_rating\": 2, \"startup\": {\"ai_nation_count\": 3, \"start_up_cash\": 10, \"ai_start_up_cash\": 20, \"ai_aggressiveness\": 3, \"start_up_independent_town\": 15, \"start_up_raw_site\": 55, \"difficulty_level\": 4 }}"
	//args = append(args, "invoke")
	args = append(args, "update_config")
	args = append(args, s2)

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in update config ...")

	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, setup.ChainCodeID, setup.ChainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	// Register for commit event
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	txResponse, err := fcutil.CreateAndSendTransaction(setup.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	fmt.Println(txResponse)
	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	return txID, nil
}

// getEventHub initilizes the event hub
func getEventHub(client fabricClient.Client) (events.EventHub, error) {
	Logger.Infof("in getEventHub()")
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new event hub: %v", err)
	}
	Logger.Infof("eventHub: %+v", eventHub)
	foundEventHub := false
	peerConfig, err := config.GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading peer config: %v", err)
	}
	Logger.Infof("peerConfig: %+v", peerConfig)
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("******* EventHub connect to peer (%s:%d) *******\n", p.EventHost, p.EventPort)
			Logger.Infof("Connecting to URL (%s:%d)\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		Logger.Infof("NO EventHub configuration found")
		return nil, fmt.Errorf("No EventHub configuration found")
	}

	Logger.Infof("FOUND eventHub: %+v", eventHub)
	return eventHub, nil
}
