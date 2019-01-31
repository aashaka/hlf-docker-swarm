package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
//	"math"
	"math/rand"
	"strconv"
	"sync"
)
var load = 1000
var numThreads = 48
var wg sync.WaitGroup

type pair struct {
	key string
	value string
}

func worker(id int, jobs <-chan pair, results chan<- int, args []string, setup *FabricSetup) {
	wg.Add(1)
	defer wg.Done()
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j.key)
		_, err := setup.clients[id].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(j.key), []byte(j.value)}})
		fmt.Println("worker", id, "finished job", j)
		if err != nil {
			results <- 0
		} else {
			results <- 1
		}
	}
}

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
	var flag int
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for k:=0; k<load; {
		flag = 1
		num := strconv.Itoa(r1.Intn(10000000))
		num2 := strconv.Itoa(r1.Intn(100))
		for _, ele := range randNum1 {
			if ele == num {
				flag=0
			}
		}
		if flag == 1 {
			randNum1 = append(randNum1,num)
			k++
			randNum2 = append(randNum2,num2)
		}
	}
	jobs := make(chan pair, load)
	results := make(chan int, load)
	for ki:=0; ki<load; ki++ {
		jobs <- pair{randNum1[ki],randNum2[ki]}
	}

	start := time.Now()

	for w := 0; w < numThreads ; w++ {
		go worker(w, jobs, results, args, setup)
	}

	close(jobs)

	wg.Wait()
	t1 := time.Now()
        elapsed1 := t1.Sub(start)
        fmt.Printf("Time elapsed1: %v\n",elapsed1)
	success:=0
	var res int
	for l := 1; l <= load; l++ {
		res = <-results
		success = success + res
	}

        t2 := time.Now()
        elapsed2 := t2.Sub(start)
        fmt.Printf("Time elapsed2: %v\n",elapsed2)
        fmt.Printf(" Able: %d\n Unable: %d\n",success, load-success)
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
