package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"io/ioutil"
	"time"
//	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
)
var load = 50000
var tload = 10000
var numThreads = 48
var wg sync.WaitGroup

type pair struct {
	key string
	value string
}

type triplet struct {
	to string
	from string
	value string
}

var randNum1 []string
var randNum2 []string

func worker(id int, jobs <-chan pair, results chan<- int, args []string, setup *FabricSetup) {
	wg.Add(1)
	defer wg.Done()
	for j := range jobs {
		go func(k pair){
			fmt.Println("open worker", id, "started  job", k.key, " at ", time.Now())
			_, err := setup.clients[id].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(k.key), []byte(k.value)}})
			fmt.Println("worker", id, "finished job", j)
			if err != nil {
				results <- 0
			} else {
				results <- 1
			}
		}(j)
		time.Sleep(100 * time.Millisecond)
	}
}


func tworker(id int, tjobs <-chan triplet, tresults chan<- int, args []string, setup *FabricSetup) {
        wg.Add(1)
        defer wg.Done()
        for j := range tjobs {
                go func(k triplet){
                        fmt.Println("transfer worker", id, "started  job", k.from, " -> ", k.to, " at ", time.Now())
                        _, err := setup.clients[id].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(k.from), []byte(k.to), []byte(k.value)}})
                        fmt.Println("worker", id, "finished job", j)
                        if err != nil {
                                tresults <- 0
                        } else {
                                tresults <- 1
                        }
                }(j)
                time.Sleep(100 * time.Millisecond)
        }
}

func (setup *FabricSetup) OpenAccounts(startidx int, length int) (string, error) {
        fmt.Println("Invoked Open")
        // Prepare arguments
        var args []string
        args = append(args, "invoke")
        args = append(args, "open")

        jsonFile, err := os.Open("state.json")
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Successfully Opened state.json")
        defer jsonFile.Close()

        var acs []string

        var result map[string]interface{}
        byteValue, _ := ioutil.ReadAll(jsonFile)
        json.Unmarshal([]byte(byteValue), &result)

        accounts := result["accounts"].(map[string]interface{})
        for key, _ := range accounts {

                // Each value is an interface{} type, that is type asserted as a string
                acs=append(acs,key)

        }

        jobs := make(chan pair, length)
        results := make(chan int, length)
        for ki:=startidx; ki<startidx+length; ki++ {
                jobs <- pair{acs[ki],"1000"}
        }

        start := time.Now()
	fmt.Println("", time.Now(), ": Starting account creation")
        for w := 0; w < numThreads ; w++ {
                go worker(w, jobs, results, args, setup)
        }
	wg.Wait()
        close(jobs)
	fmt.Println("", time.Now(), ": Done opening accounts")

        t1 := time.Now()
        elapsed1 := t1.Sub(start)
        fmt.Printf("Time elapsed1: %v\n",elapsed1)
        success:=0
        var res int
        for l := 1; l <= length; l++ {
                res = <-results
                success = success + res
        }

        t2 := time.Now()
        elapsed2 := t2.Sub(start)
        fmt.Printf("Time elapsed2: %v\n",elapsed2)
        fmt.Printf(" Able: %d\n Unable: %d\n",success, load-success)



	return "", nil
}



func (setup *FabricSetup) TransferAccounts(startidx int, length int) (string, error) {
        fmt.Println("Invoked Open")
        // Prepare arguments
        var args []string
        args = append(args, "invoke")
        args = append(args, "transfer")



        jsonFile, err := os.Open("state.json")
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Successfully Opened state.json")
        defer jsonFile.Close()

        var acs []string

        var result map[string]interface{}
        byteValue, _ := ioutil.ReadAll(jsonFile)
        json.Unmarshal([]byte(byteValue), &result)

        accounts := result["accounts"].(map[string]interface{})
        for key, _ := range accounts {

                // Each value is an interface{} type, that is type asserted as a string
                acs=append(acs,key)

        }


	var from_el []string
	var to_el []string

	for k:=startidx; k<length*2; k=k+2 {
		from_el = append(from_el, acs[k])
		to_el = append(to_el,acs[k+1])
	}


	tjobs := make(chan triplet, length)
	tresults := make(chan int, length)
	for ki:=0; ki<length; ki++ {
		tjobs <- triplet{from_el[ki],to_el[ki],"100"}
	}

	start := time.Now()

	fmt.Println(time.Now(), ": Starting transfer tx")
	for w := 0; w < numThreads ; w++ {
		go tworker(w, tjobs, tresults, args, setup)
	}

	close(tjobs)

	wg.Wait()
	fmt.Println(time.Now(), ": Transfers finished")
        t1 := time.Now()
        elapsed1 := t1.Sub(start)
        fmt.Printf("Time elapsed1: %v\n",elapsed1)
        success:=0
        var res int
        for l := 1; l <= length; l++ {
                res = <-tresults
                success = success + res
        }

        t2 := time.Now()
        elapsed2 := t2.Sub(start)
        fmt.Printf("Time elapsed2: %v\n",elapsed2)
        fmt.Printf(" Able: %d\n Unable: %d\n",success, tload-success)

	return "", nil
}

// InvokeOpen
func (setup *FabricSetup) InvokeOpen(account string, value string) (string, error) {

	fmt.Println("Invoked Open")
	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "open")
	args = append(args, account)
	args = append(args, value)

//	eventID := "eventOpen"

	// Add data that will be visible in the proposal, like a description of the invoke request

// 	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
// 	if err != nil {
// 		return "", err
// 	}
//	defer setup.event.Unregister(reg)
//	var randNum1 []string
//	var randNum2 []string
	var flag int
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for k:=0; k<load; {
		flag = 1
		num := strconv.Itoa(r1.Intn(10000000))
		num2 := strconv.Itoa(10000 + r1.Intn(100))
//		fmt.Printf("num1: %d, k: %d\n", num, k)
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
	fmt.Println("Load is ",load)
	jobs := make(chan pair, load)
	results := make(chan int, load)
	for ki:=0; ki<load; ki++ {
//		fmt.Println("Got key", ki)
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

//	eventID := "eventTransfer"
//
//	// Add data that will be visible in the proposal, like a description of the invoke request
//	transientDataMap := make(map[string][]byte)
//	transientDataMap["result"] = []byte("Transient data in transfer")
//
//	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
//	if err != nil {
//		return "", err
//	}
//	defer setup.event.Unregister(reg)

	var from_el []string
	var to_el []string
	//s1 := rand.NewSource(time.Now().UnixNano())
        //r1 := rand.New(s1)
	for k:=0; k<tload*2; k=k+2 {
		from_el = append(from_el, randNum1[k])
		to_el = append(to_el,randNum1[k+1])
	}


	tjobs := make(chan triplet, tload)
	tresults := make(chan int, tload)
	for ki:=0; ki<tload; ki++ {
		tjobs <- triplet{from_el[ki],to_el[ki],"100"}
	}

	start := time.Now()

	for w := 0; w < numThreads ; w++ {
		go tworker(w, tjobs, tresults, args, setup)
	}

	close(tjobs)

	wg.Wait()
	t1 := time.Now()
        elapsed1 := t1.Sub(start)
        fmt.Printf("Time elapsed1: %v\n",elapsed1)
	success:=0
	var res int
	for l := 1; l <= tload; l++ {
		res = <-tresults
		success = success + res
	}

        t2 := time.Now()
        elapsed2 := t2.Sub(start)
        fmt.Printf("Time elapsed2: %v\n",elapsed2)
        fmt.Printf(" Able: %d\n Unable: %d\n",success, tload-success)


// Create a request (proposal) and send it
//	response, err := setup.clients[0].Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])}, TransientMap: transientDataMap})
//	if err != nil {
//		return "", fmt.Errorf("failed to open account: %v", err)
//	}

	// Wait for the result of the submission
//	select {
//	case ccEvent := <-notifier:
//		fmt.Printf("Received CC event: %s\n", ccEvent)
//	case <-time.After(time.Second * 20):
//		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
//	}

	return "", nil
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
