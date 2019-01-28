package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
	"math"
	"math/rand"
	"strconv"
	"sync"
)
var load = 200
var numThreads = 48
// InvokeOpen
func (setup *FabricSetup) InvokeOpen(account string, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "open")
	args = append(args, account)
	args = append(args, value)

//	eventID := "eventOpen"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in open")

// 	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
// 	if err != nil {
// 		return "", err
// 	}
//	defer setup.event.Unregister(reg)
	var randNum1 []string
	var randNum2 []string
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for k:=0; k<load; k++ {
		num := strconv.Itoa(r1.Intn(10000000))
		num2 := strconv.Itoa(r1.Intn(100))
		randNum1 = append(randNum1,num)
		randNum2 = append(randNum2,num2)
	}
///	resch := make(chan channel.Response, load+4)
///	errch := make(chan error, load+4)
	var able int
	var unable int
	var wg sync.WaitGroup
	start := time.Now()
///	ablech := make(chan int, load+4)
///	unablech := make(chan int, load+4)
	chunksz := load/numThreads
	tx := 0
	for i:=0; i<numThreads; i++ {
		wg.Add(1)
		j := i
		go func() {
			defer wg.Done()
			for tx=j*chunksz; float64(tx)<math.Min(float64((j+1)*chunksz),float64(load)); tx++ {
				// Create a request (proposal) and send it
				fmt.Printf("Sending transaction %d via client %d\n",tx,j)
//				fmt.Printf("%s: %s",randNum1[j], randNum2[j])
				_, err := setup.clients[j].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(randNum1[tx]), []byte(randNum2[tx])}, TransientMap: transientDataMap})
				if err != nil {
///					errch<-fmt.Errorf("failed to open account: %v", err)
///					unablech<-j
				} else {
	//				fmt.Printf("TXID: %v", response.TransactionID)
///					resch<-response
///					ablech<-j
				}
			}
/*			if j == numThreads-1 {
				for ; tx<load; tx++ {
					// Create a request (proposal) and send it
	                                fmt.Printf("Sending transaction %d via client %d",tx,j)
	//                              fmt.Printf("%s: %s",randNum1[j], randNum2[j])
	                                _, err := setup.clients[j].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(randNum1[tx]), []byte(randNum2[tx])}, TransientMap: transientDataMap})
	                                if err != nil {
///	                                        errch<-fmt.Errorf("failed to open account: %v", err)
///	                                        unablech<-j
	                                } else {
	        //                              fmt.Printf("TXID: %v", response.TransactionID)
///	                                        resch<-response
///	                                        ablech<-j
	                                }
				}
			}
*/
		}()
	}
	wg.Wait()
	t1 := time.Now()
        elapsed1 := t1.Sub(start)
        fmt.Printf("Time elapsed1: %v\n",elapsed1)
/*	for l := 1; l <= numThreads; l++ {
		select {
			case _ = <-ablech:
				able = able + 1
		//		fmt.Printf("######## Received response for transaction: %d\n", ind)
			case _ = <-unablech:
				unable = unable + 1
		//		fmt.Printf("######## No response for transaction: %d\n", indn)
		}
	}
*/
        t2 := time.Now()
        elapsed2 := t2.Sub(start)
        fmt.Printf("Time elapsed2: %v\n",elapsed2)
        fmt.Printf(" Able: %d\n Unable: %d\n",able, unable)
// 	// Wait for the result of the submission
// 	select {
// 	case ccEvent := <-notifier:
// 		fmt.Printf("Received CC event: %s\n", ccEvent)
// 	case <-time.After(time.Second * 20):
// 		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
// 	}

	return "", nil
}

// InvokeQuery
func (setup *FabricSetup) InvokeQuery(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, value)

	// Create a request (proposal) and send it
	response, err := setup.clients[0].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])} })
	if err != nil {
		return "", fmt.Errorf("failed to open account: %v", err)
	}

	return string(response.Payload), nil
}

// InvokeOpen
func (setup *FabricSetup) InvokeTransfer(sender string, receiver string, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "transfer")
	args = append(args, sender)
	args = append(args, receiver)
	args = append(args, value)

	eventID := "eventTransfer"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in transfer")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.clients[0].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to open account: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}

// InvokeDelete
func (setup *FabricSetup) InvokeDelete(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "delete")
	args = append(args, value)

	eventID := "eventDelete"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in delete")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.clients[0].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to open account: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}
