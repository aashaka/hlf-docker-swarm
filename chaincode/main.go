package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const ERROR_SYSTEM = "{\"code\":300, \"reason\": \"system error: %s\"}"
const ERROR_WRONG_FORMAT = "{\"code\":301, \"reason\": \"command format is wrong\"}"
const ERROR_ACCOUNT_EXISTING = "{\"code\":302, \"reason\": \"account already exists\"}"
const ERROR_ACCOUNT_ABNORMAL = "{\"code\":303, \"reason\": \"abnormal account\"}"
const ERROR_MONEY_NOT_ENOUGH = "{\"code\":304, \"reason\": \"account's money is not enough\"}"



type SimpleChaincode struct {

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// nothing to do
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if args[0] == "open" {
		return t.Open(stub, args)
	}
	if args[0] == "delete" {
		return t.Delete(stub, args)
	}
	if args[0] == "query" {
		return t.Query(stub, args)
	}
	if args[0] == "transfer" {
		return t.Transfer(stub, args)
	}

	return shim.Error(ERROR_WRONG_FORMAT)
}

// open an account, should be [open account money]
func (t *SimpleChaincode) Open(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	account  := args[1]
	money,err := stub.GetState(account)
	if money != nil {
		return shim.Error(ERROR_ACCOUNT_EXISTING)
	}

	_,err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	err = stub.PutState(account, []byte(args[2]))
	if err != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}

	err = stub.SetEvent("eventOpen", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// delete an account, should be [delete account]
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	err := stub.DelState(args[1])
	if err != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}

	err = stub.SetEvent("eventDelete", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query current money of the account,should be [query accout]
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	money, err := stub.GetState(args[1])
	if err != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}

	if money == nil {
		return shim.Error(ERROR_ACCOUNT_ABNORMAL)
	}

	return shim.Success(money)
}

// transfer money from account1 to account2, should be [transfer account1 account2 money]
func (t *SimpleChaincode) Transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error(ERROR_WRONG_FORMAT)
	}
	money, err := strconv.Atoi(args[3])

	if err != nil {
		return shim.Error(ERROR_WRONG_FORMAT)
	}

	moneyBytes1, err1 := stub.GetState(args[1])
	moneyBytes2, err2 := stub.GetState(args[2])
	if err1 != nil || err2 != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}
	if moneyBytes1 == nil || moneyBytes2 == nil {
		return shim.Error(ERROR_ACCOUNT_ABNORMAL)
	}

	money1, _ := strconv.Atoi(string(moneyBytes1))
	money2, _ := strconv.Atoi(string(moneyBytes2))
	if money1 < money {
		return shim.Error(ERROR_MONEY_NOT_ENOUGH)
	}

	money1 -= money
	money2 += money

	err = stub.PutState(args[1], []byte(strconv.Itoa(money1)))
	if err != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}

	err = stub.PutState(args[2], []byte(strconv.Itoa(money2)))
	if err != nil {
		s := fmt.Sprintf(ERROR_SYSTEM, err.Error())
		return shim.Error(s)
	}

	err = stub.SetEvent("eventTransfer", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func  main()  {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %v \n", err)
	}

}


// // HeroesServiceChaincode implementation of Chaincode
// type HeroesServiceChaincode struct {
// }

// // Init of the chaincode
// // This function is called only one when the chaincode is instantiated.
// // So the goal is to prepare the ledger to handle future requests.
// func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
// 	fmt.Println("########### HeroesServiceChaincode Init ###########")

// 	// Get the function and arguments from the request
// 	function, _ := stub.GetFunctionAndParameters()

// 	// Check if the request is the init function
// 	if function != "init" {
// 		return shim.Error("Unknown function call")
// 	}

// 	// Put in the ledger the key/value hello/world
// 	err := stub.PutState("hello", []byte("world"))
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	// Return a successful message
// 	return shim.Success(nil)
// }

// // Invoke
// // All future requests named invoke will arrive here.
// func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
// 	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

// 	// Get the function and arguments from the request
// 	function, args := stub.GetFunctionAndParameters()

// 	// Check whether it is an invoke request
// 	if function != "invoke" {
// 		return shim.Error("Unknown function call")
// 	}

// 	// Check whether the number of arguments is sufficient
// 	if len(args) < 1 {
// 		return shim.Error("The number of arguments is insufficient.")
// 	}

// 	// In order to manage multiple type of request, we will check the first argument.
// 	// Here we have one possible argument: query (every query request will read in the ledger without modification)
// 	if args[0] == "query" {
// 		return t.query(stub, args)
// 	}

// 	// The update argument will manage all update in the ledger
// 	if args[0] == "invoke" {
// 		return t.invoke(stub, args)
// 	}

// 	// If the arguments given don’t match any function, we return an error
// 	return shim.Error("Unknown action, check the first argument")
// }

// // query
// // Every readonly functions in the ledger will be here
// func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	fmt.Println("########### HeroesServiceChaincode query ###########")

// 	// Check whether the number of arguments is sufficient
// 	if len(args) < 2 {
// 		return shim.Error("The number of arguments is insufficient.")
// 	}

// 	// Like the Invoke function, we manage multiple type of query requests with the second argument.
// 	// We also have only one possible argument: hello
// 	if args[1] == "hello" {

// 		// Get the state of the value matching the key hello in the ledger
// 		state, err := stub.GetState("hello")
// 		if err != nil {
// 			return shim.Error("Failed to get state of hello")
// 		}

// 		// Return this value in response
// 		return shim.Success(state)
// 	}

// 	// If the arguments given don’t match any function, we return an error
// 	return shim.Error("Unknown query action, check the second argument.")
// }

// // invoke
// // Every functions that read and write in the ledger will be here
// func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	fmt.Println("########### HeroesServiceChaincode invoke ###########")

// 	if len(args) < 2 {
// 		return shim.Error("The number of arguments is insufficient.")
// 	}

// 	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
// 	if args[1] == "hello" && len(args) == 3 {

// 		// Write the new value in the ledger
// 		err := stub.PutState("hello", []byte(args[2]))
// 		if err != nil {
// 			return shim.Error("Failed to update state of hello")
// 		}

// 		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
// 		err = stub.SetEvent("eventInvoke", []byte{})
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		// Return this value in response
// 		return shim.Success(nil)
// 	}

// 	// If the arguments given don’t match any function, we return an error
// 	return shim.Error("Unknown invoke action, check the second argument.")
// }

// func main() {
// 	// Start the chaincode and make it ready for futures requests
// 	err := shim.Start(new(HeroesServiceChaincode))
// 	if err != nil {
// 		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
// 	}
// }