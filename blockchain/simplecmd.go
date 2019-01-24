package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
	"math/rand"
	"strconv"
	"sync"
)


// InvokeOpen
func (setup *FabricSetup) WriteRand() (time.Duration, error) {
//args: [, '1', '0'.repeat(5), '9'.repeat(9), '8', '20', '','']
//	s1 := rand.NewSource(time.Now().UnixNano())
//	r1 := rand.New(s1)
//	num := strconv.Itoa(r1.Intn(100))
	// Prepare arguments
	var args []string
	args = append(args, "WriteRandom")
	args = append(args, "num")
	args = append(args, "1")
	args = append(args, "00000")
	args = append(args, "999999999")
	args = append(args, "8")
	args = append(args, "20")
	args = append(args, "")
	args = append(args, "")

	eventID := "eventWR"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in open")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return time.Duration(0), err
	}
	defer setup.event.Unregister(reg)

	start := time.Now()
	var wg sync.WaitGroup
	c1 := make(chan channel.Response, 100)
	errch := make(chan error)
	defer close(c1)
	for i := 1; i <= 100; i++ {
		j := i
		go func() {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			num := strconv.Itoa(r1.Intn(100))
			wg.Add(1)
			defer wg.Done()
			fmt.Println("########## j is ", j)
			response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(num), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7]), []byte(args[8])}, TransientMap: transientDataMap})
			if err != nil {
				errch<-fmt.Errorf("failed to randomly write: %v",err)
			} else {
				c1<-response
			}
		}()
	}
	for j := 1; j <= 100; j++ {
		select {
			case res := <-c1:
				fmt.Printf("######## Received response for transaction: %s, %d more to go\n", res.TransactionID, 100-j)
			case err := <-errch:
				fmt.Printf("####### Received error %s, %d more to go", err, 100-j)
		}
	}


	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Time passed for 100 random write transactions is:", elapsed)
	// Create a request (proposal) and send it
	// response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7]), []byte(args[8])}, TransientMap: transientDataMap})
	// if err != nil {
	// 	return "", fmt.Errorf("failed to open account: %v", err)
	// }

	// Wait for the result of the submission
	for k := 1; k <= 100; k++ {
		select {
		case ccEvent := <-notifier:
			fmt.Printf("Received CC event: %s\n", ccEvent)
		case <-time.After(time.Second * 20):
			return time.Duration(0), fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		}
	}

	return elapsed, nil
}



// InvokeOpen
func (setup *FabricSetup) InvokeOpen(account string, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "open")
	args = append(args, account)
	args = append(args, value)

	eventID := "eventOpen"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in open")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
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

// InvokeQuery
func (setup *FabricSetup) InvokeQuery(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, value)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])} })
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
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}, TransientMap: transientDataMap})
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
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}, TransientMap: transientDataMap})
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
