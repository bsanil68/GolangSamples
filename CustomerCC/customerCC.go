package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Customer struct {
	CustomerID       string `json:"CustomerID"`
	LoanID           string `json:"LoanID"`
	BorrowerName     string `json:"BorrowerName"`
	BorrowerSODOWO   string `json:"BorrowerSODOWO"`
	AgeOfBorrower    string `json:"AgeOfBorrower"`
	CoBorrowerName   string `json:"CoBorrowerName"`
	CoBorrowerSODOWO string `json:"CoBorrowerSODOWO"`
	AgeOfCoBorrower  string `json:"AgeOfCoBorrower"`
	CustomerHash     string `json:"CustomerHash"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *Customer) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *Customer) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "CreateCustomer" {

		result, err = CreateCustomer(stub, args)
	}
	if fn == "GetCustomerByCustomerID" {

		result, err = GetCustomerByCustomerID(stub, args)
	}
	if fn == "GetCustomerByCustomerHash" {

		result, err = GetCustomerByCustomerHash(stub, args)
	}
	if fn == "GetCustomerByLoanID" {

		result, err = GetCustomerByLoanID(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

/**
         The functions   gets  customers  stored in ledger  based  on  id


**/
func GetCustomerByCustomerID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetCustomerFromID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing CustomerID")
	}

	//form client customer id is   the key  which is   created
	//  by a  custom logic
	var customerID = args[0]

	value, err := stub.GetState(customerID)
	if err != nil {
		fmt.Println("Could not get customerid  with id "+customerID+" from ledger", err)
		return "", errors.New("Missing CustomerID")
	}

	return string(value), nil
}

func GetCustomerByCustomerHash(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetCustomerByCustomerHash")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing CustomerHash argument")
	}

	var customer_Hash = args[0]

	getCustomerByCustomerHashItr, err := stub.GetStateByPartialCompositeKey("customerhashIndex", []string{customer_Hash})

	if err != nil {
		return "", fmt.Errorf("Could not get composite key for this customer_Hash")
	}
	defer getCustomerByCustomerHashItr.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("{\"CustomerForCustomerHash\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getCustomerByCustomerHashItr.HasNext() {
		queryResponse, err := getCustomerByCustomerHashItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next loan Data")
		}

		// Add a comma before arrat members, supress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the customerId from the customerId~loanId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite key parts")
		}
		returnedCustomerHash := compositeKeyParts[0]
		returnedCustomerID := compositeKeyParts[1]

		fmt.Printf("- found a from index:%s customerHash:%s customerID:%s", objectType, returnedCustomerHash, returnedCustomerID)

		value, err := stub.GetState(returnedCustomerID)
		if err != nil {
			fmt.Println("Couldn't get customer for Id"+returnedCustomerID+"from ledger", err)
			return "", errors.New("Missing customerID in ledger")
		}

		// Returned Record is a JSON object. So write as is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetCustomerByLoanID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetCustomerByLoanID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing LoanID argument")
	}

	var loanId = args[0]

	getCustomerByLoanIDItr, err := stub.GetStateByPartialCompositeKey("customerloanIndex", []string{loanId})

	if err != nil {
		return "", fmt.Errorf("Could not get composite key for this loanId")
	}
	defer getCustomerByLoanIDItr.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("{\"CustomerForLoan\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getCustomerByLoanIDItr.HasNext() {
		queryResponse, err := getCustomerByLoanIDItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next loan Data")
		}

		// Add a comma before arrat members, supress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the customerId from the customerId~loanId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite key parts")
		}
		returnedloanID := compositeKeyParts[0]
		returnedCustomerID := compositeKeyParts[1]

		fmt.Printf("- found a from index:%s loanID:%s customerID:%s", objectType, returnedloanID, returnedCustomerID)

		value, err := stub.GetState(returnedCustomerID)
		if err != nil {
			fmt.Println("Couldn't get customer for id "+returnedCustomerID+"from ledger", err)
			return "", errors.New("Missing customerID in ledger")
		}

		// Returned Record is a JSON object. So write as is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func CreateCustomer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CreateCustomer")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	// json args input musts match  the case and  spelling exactly
	// get all the  arguments
	var custCid = args[0]
	fmt.Println("the customerid is" + custCid)
	var loanId = args[1]
	fmt.Println("The LoanId is" + loanId)
	var borrowerName = args[2]
	fmt.Println("the borrowerName is" + borrowerName)
	var borrowerSODOWO = args[3]
	fmt.Println("the borrowerSODOWO is" + borrowerSODOWO)
	var ageOfBorrower = args[4]
	fmt.Println("the ageOfBorrower is" + ageOfBorrower)
	var coBorrowerName = args[5]
	fmt.Println("the coBorrowerName is" + coBorrowerName)
	var coBorrowerSODOWO = args[6]
	fmt.Println("the coBorrowerSODOWO is" + coBorrowerSODOWO)
	var ageOfCoBorrower = args[7]
	fmt.Println("the ageOfCoBorrower is" + ageOfCoBorrower)
	var customerHash = args[8]
	fmt.Println("the customerHash is" + customerHash)

	//assigning to struct the variables
	customerStruct := Customer{
		CustomerID:       custCid,
		LoanID:           loanId,
		BorrowerName:     borrowerName,
		BorrowerSODOWO:   borrowerSODOWO,
		AgeOfBorrower:    ageOfBorrower,
		CoBorrowerName:   coBorrowerName,
		CoBorrowerSODOWO: coBorrowerSODOWO,
		AgeOfCoBorrower:  ageOfCoBorrower,
		CustomerHash:     customerHash,
	}

	fmt.Println("the struct values are custid  " + customerStruct.CustomerID + "loanID" + customerStruct.LoanID + "borrowername" + customerStruct.BorrowerName + "age" + customerStruct.AgeOfBorrower)

	custStructBytes, err := json.Marshal(customerStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	custerr := stub.PutState(custCid, []byte(custStructBytes))
	if custerr != nil {
		fmt.Println("Could not save customer data to ledger", custerr)
		return "", fmt.Errorf("Could not save customer data to ledger")
	}

	compIndexName := "customerloanIndex"
	customerLoanCompKey, err := stub.CreateCompositeKey(compIndexName, []string{customerStruct.LoanID, customerStruct.CustomerID})
	if err != nil {
		fmt.Println("Could not index composite key for customerId~loanId", err)
		return "", fmt.Errorf("Could not index composite key for customerId~loanId map")
	}

	value := []byte{0x00}
	stub.PutState(customerLoanCompKey, value)

	indexName2 := "customerhashIndex"
	customerhashIndexKey, err := stub.CreateCompositeKey(indexName2, []string{customerStruct.CustomerHash, customerStruct.CustomerID})
	if err != nil {
		fmt.Println("Could not index composite key for customerId~CustomerHash", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value2 := []byte{0x00}
	stub.PutState(customerhashIndexKey, value2)

	var customEvent = "{eventType: 'CreateCustomer', description: CustomerHash: " + customerStruct.CustomerHash + "CustomerID" + customerStruct.CustomerID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved customer")
	return args[0], nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(Customer))
	{
		fmt.Printf("Error starting enrollcustomer chaincode: %s", err)
	}
}
