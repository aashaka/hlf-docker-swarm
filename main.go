package main

import (
	"fmt"
	"github.com/hyperledger/hlf-docker-swarm/blockchain"
	// "github.com/hyperledger/hlf-docker-swarm/web"
	// "github.com/hyperledger/hlf-docker-swarm/web/controllers"
	"os"
//	"time"
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
		ChaincodePath:   "github.com/hyperledger/hlf-docker-swarm/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",
		// User parameters
		UserName: "User1",
	}

	if(os.Args[1] == "0") {
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

	// 	response, err := fSetup.QueryHello()
	// 	if err != nil {
	// 		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	// 	} else {
	// 		fmt.Printf("Response from the query hello: %s\n", response)
	// 	}
	// 	var unable int
	// 	var able int
	// 	start := time.Now()
	// 	for i := 1; i <= 100; i++ {
	// 		j := i
	// 		go func() {
				// Query the chaincode
//				txId, err := fSetup.InvokeOpen("Alice", "100")
//				if err != nil {
//					fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	//				unable = unable + 1
//				} else {
//					fmt.Printf("Response from the invoke open, transaction ID: %s\n", txId)
	//				able = able + 1
//				}
	//		}()
	//	}

	// 	t := time.Now()
	// 	elapsed := t.Sub(start)
	// 	fmt.Printf("Time elapsed: %v\n",elapsed)
	// 	fmt.Printf(" Able: %d\n Unable: %d\n",able, unable)
	//	txId, err = fSetup.InvokeOpen("Bob", "100")
	//	if err != nil {
	//		fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	//	} else {
	//		fmt.Printf("Response from the invoke open, transaction ID: %s\n", txId)
	//	}
	//
	//	response, err := fSetup.InvokeQuery("Alice")
	//	if err != nil {
	//		fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	//	} else {
	//		fmt.Printf("Response from the query Alice: %s\n", response)
	//	}
	//
	//	response, err = fSetup.InvokeQuery("Bob")
	//	if err != nil {
	//		fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	//	} else {
	//		fmt.Printf("Response from the query Bob: %s\n", response)
	//	}
	//
	//	txId, err = fSetup.InvokeTransfer("Bob", "Alice", "20")
	//	if err != nil {
	//		fmt.Printf("Unable to invoke transfer on the chaincode: %v\n", err)
	//	} else {
	//		fmt.Printf("Successfully invoked transfer, transaction ID: %s\n", txId)
	//	}
	//
	//	txId, err = fSetup.InvokeDelete("Bob")
	//	if err != nil {
	//		fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	//	} else {
	//		fmt.Printf("Response from the delete, transaction ID: %s\n", txId)
	//	}

		// // Launch the web application listening
		// app := &controllers.Application{
		// 	Fabric: &fSetup,
		// }
		// web.Serve(app)
	} else if (os.Args[1]=="1") {
				err := fSetup.Continue()
                                txId, err := fSetup.InvokeOpen("Alice", "100")
                                if err != nil {
                                        fmt.Printf("Unable to open account on the chaincode: %v\n", err)
        //                              unable = unable + 1
                                } else {
                                        fmt.Printf("Response from the invoke open, transaction ID: %s\n", txId)
        //                              able = able + 1
                                }
                                txId, err = fSetup.InvokeTransfer("Alice", "Bob","100")
                                if err != nil {
                                        fmt.Printf("Unable to transfer on the chaincode: %v\n", err)
        //                              unable = unable + 1
                                } else {
                                        fmt.Printf("Response from the invoke transfer, transaction ID: %s\n", txId)
        //                              able = able + 1
                                }
	} else {
		fmt.Println("0 or 1 for now")
	}
}
