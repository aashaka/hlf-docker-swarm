package main

import (
	"fmt"
	"github.com/hyperledger/hlf-docker-swarm/blockchain"
	// "github.com/hyperledger/hlf-docker-swarm/web"
	// "github.com/hyperledger/hlf-docker-swarm/web/controllers"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer0.example.com",

		// Channel parameters
		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/hlf-docker-swarm/network/config/mychannel.tx",
		// AnchorPeersConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/hlf-docker-swarm/network/config/Org1MSPanchors_mychannel.tx",
		// Chaincode parameters
		ChainCodeID:     "hlf-docker-swarm",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hyperledger/hlf-docker-swarm/chaincodes/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

		// Query the chaincode
	response, err := fSetup.InvokeOpen("Alice", "100")
	if err != nil {
		fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the invoke open: %s\n", response)
	}

	response, err := fSetup.InvokeOpen("Bob", "100")
	if err != nil {
		fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the invoke open: %s\n", response)
	}

	// Invoke the chaincode
	txId, err := fSetup.InvokeTransfer("Bob", "Alice", "20")
	if err != nil {
		fmt.Printf("Unable to invoke transfer on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Successfully invoked transfer, transaction ID: %s\n", txId)
	}

	// Query again the chaincode
	response, err = fSetup.InvokeDelete('Bob')
	if err != nil {
		fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the delete: %s\n", response)
	}

	// // Launch the web application listening
	// app := &controllers.Application{
	// 	Fabric: &fSetup,
	// }
	// web.Serve(app)
}