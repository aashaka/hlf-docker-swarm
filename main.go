package main

import (
	"fmt"
	"github.com/hyperledger/hlf-docker-swarm/blockchain"
	// "github.com/hyperledger/hlf-docker-swarm/web"
	// "github.com/hyperledger/hlf-docker-swarm/web/controllers"
	"os"
)

//func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
//  sc := bufio.NewScanner(r)
//  for sc.Scan() {
//      lastLine++
//      if lastLine == lineNum {
//          // you can return sc.Bytes() if you need output in []bytes
//          return sc.Text(), lastLine, sc.Err()
//      }
//  }
//  return line, lastLine, io.EOF
//}

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
		ChaincodePath:   "github.com/hyperledger/hlf-docker-swarm/chaincode",
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
	// for i := 0; i < 100; i++ {
	// 	num := i
	// 	go func() {
	// 		txId, err := fSetup.WriteRandom(prng().toString().substr(2,4), "1", "0".repeat(5), "9".repeat(9), "8", "20", "","")
	// 		if err != nil {
	// 			fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	// 		} else {
	// 			fmt.Printf("Response from the invoke open, transaction ID: %s\n", txId)
	// 		}
	// 	}
	// }	

	elapsed, err := fSetup.WriteRand()
	if err != nil {
		fmt.Printf("Unable to open account on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the write rand received, time taken = %v", elapsed)
	}

	// response, err := fSetup.InvokeQuery("Alice")
	// if err != nil {
	// 	fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Response from the query Alice: %s\n", response)
	// }

	// response, err = fSetup.InvokeQuery("Bob")
	// if err != nil {
	// 	fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Response from the query Bob: %s\n", response)
	// }

	// txId, err = fSetup.InvokeTransfer("Bob", "Alice", "20")
	// if err != nil {
	// 	fmt.Printf("Unable to invoke transfer on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Successfully invoked transfer, transaction ID: %s\n", txId)
	// }

	// txId, err = fSetup.InvokeDelete("Bob")
	// if err != nil {
	// 	fmt.Printf("Unable to delete account on the chaincode: %v\n", err)
	// } else {
	// 	fmt.Printf("Response from the delete, transaction ID: %s\n", txId)
	// }

	// // Launch the web application listening
	// app := &controllers.Application{
	// 	Fabric: &fSetup,
	// }
	// web.Serve(app)
}
