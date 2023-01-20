package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type EMISchedule struct {
	LoanID                      string `json:"LoanID"`
	EmiID                       string `json:"EmiID"`
	InstallmentNumber           string `json:"InstallmentNumber"`
	ScheduleDate                string `json:"ScheduleDate"`
	EmiAmount                   string `json:"EmiAmount"`
	Installmentprincipalamount  string `json:"Installmentprincipalamount"`
	Installmentinterestamount   string `json:"Installmentinterestamount"`
	OutstandingPrincipalBalance string `json:"OutstandingPrincipalBalance"`
	Status                      string `json:"Status"`
	Currency                    string `json:"Currency"`
	EMIHash                     string `json:"EMIHash"`
}

func (t *EMISchedule) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *EMISchedule) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "CreateEMISchedule" {

		result, err = CreateEMISchedule(stub, args)
	}
	if fn == "GetEMIScheduleByEmiID" {

		result, err = GetEMIScheduleByEmiID(stub, args)
	}
	if fn == "GetEMIScheduleByEMIHash" {

		result, err = GetEMIScheduleByEMIHash(stub, args)
	}
	if fn == "GetEMIScheduleByLoanID" {

		result, err = GetEMIScheduleByLoanID(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func GetEMIScheduleByEmiID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetEMIScheduleByEmiID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing EMIScheduleID")
	}

	var emiID = args[0]

	value, err := stub.GetState(emiID)
	if err != nil {
		fmt.Println("Could not get emiID  with id "+emiID+" from ledger", err)
		return "", errors.New("Missing emiID")
	}

	return string(value), nil
}

func GetEMIScheduleByEMIHash(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetEMIScheduleByEMIHash")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing EMIHash")
	}

	var EMI_Hash = args[0]

	emiScheduleHashItr, err := stub.GetStateByPartialCompositeKey("EMIhashIndex", []string{EMI_Hash})
	if err != nil {
		return "", fmt.Errorf("Could not get EMI_Hash for emi schedule")
	}
	defer emiScheduleHashItr.Close()

	//  the root  for  json so that  it cane unmarshalled  to pojo

	var buffer bytes.Buffer

	buffer.WriteString("{\"EMIScheduleHash\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	// Iterate through result set and check
	for emiScheduleHashItr.HasNext() {
		// Note that we don't get the value (2nd return variable), we'll just get the loan id from the composite key
		queryResponse, err := emiScheduleHashItr.Next()
		if err != nil {
			return "", fmt.Errorf("no emi id for the hash ")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// get the investor id and month from InvestorActualPaymentIndex composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			fmt.Println("- split error   ")
			return "", fmt.Errorf("splitting error ")
		}

		returnedEMIHash := compositeKeyParts[0]
		returnedEmiID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s emiid:%s loanid:%s\n", objectType, returnedEmiID, returnedEMIHash)
		value, err := stub.GetState(returnedEmiID)
		if err != nil {
			fmt.Println("Could not get emiid for the id "+returnedEmiID+" from ledger", err)
			return "", errors.New("Missing emiid")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetEMIScheduleByLoanID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetEMIScheduleByLoanID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing LoanID")
	}

	var loanId = args[0]

	emiScheduleItr, err := stub.GetStateByPartialCompositeKey("EMIScheduleIndex", []string{loanId})
	if err != nil {
		return "", fmt.Errorf("Could not get loan id for emi schedule")
	}
	defer emiScheduleItr.Close()

	//  the root  for  json so that  it cane unmarshalled  to pojo

	var buffer bytes.Buffer

	buffer.WriteString("{\"EMIScheduleLoan\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	// Iterate through result set and check
	for emiScheduleItr.HasNext() {
		// Note that we don't get the value (2nd return variable), we'll just get the loan id from the composite key
		queryResponse, err := emiScheduleItr.Next()
		if err != nil {
			return "", fmt.Errorf("no emi id for the loan ")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// get the investor id and month from InvestorActualPaymentIndex composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			fmt.Println("- split error   ")
			return "", fmt.Errorf("splitting error ")
		}

		returnedLoanID := compositeKeyParts[0]
		returnedEmiID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s emiid:%s loanid:%s\n", objectType, returnedEmiID, returnedLoanID)
		value, err := stub.GetState(returnedEmiID)
		if err != nil {
			fmt.Println("Could not get emiid for the loan "+returnedEmiID+" from ledger", err)
			return "", errors.New("Missing emiid")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func CreateEMISchedule(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	fmt.Println("Entering CreateEMISchedule")
	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var loanID = args[0]
	fmt.Println("the loanId is" + loanID)
	var emiID = args[1]
	fmt.Println("the emiID is" + emiID)
	var installNum = args[2]
	fmt.Println("the  installmentnumber is" + installNum)
	var schDate = args[3]
	fmt.Println("the scheduledate is" + schDate)
	var emiAmt = args[4]
	fmt.Println("the EMI amount is" + emiAmt)
	var installPrinAmount = args[5]
	fmt.Println("the installment principle amount is" + installPrinAmount)
	var installIntAmt = args[6]
	fmt.Println("the installmentinterestamount is" + installIntAmt)
	var outstandPrincp = args[7]
	fmt.Println("the outstanding principle balance is" + outstandPrincp)
	var currency = args[8]
	fmt.Println("the currency is" + currency)
	var emiHash = args[9]
	fmt.Println("the emiHash is" + emiHash)

	EMIScheduleStruct := EMISchedule{
		LoanID:                      loanID,
		EmiID:                       emiID,
		InstallmentNumber:           installNum,
		ScheduleDate:                schDate,
		EmiAmount:                   emiAmt,
		Installmentprincipalamount:  installPrinAmount,
		Installmentinterestamount:   installIntAmt,
		OutstandingPrincipalBalance: outstandPrincp,
		Status:                      "active",
		Currency:                    currency,
		EMIHash:                     emiHash,
	}

	fmt.Println("the struct values are custid  " + EMIScheduleStruct.LoanID + "Emi ID" + EMIScheduleStruct.EmiID + "installmentNumber" + EMIScheduleStruct.InstallmentNumber + " scheduleDate" + EMIScheduleStruct.ScheduleDate)

	EmischdStructBytes, err := json.Marshal(EMIScheduleStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	emiScherr := stub.PutState(emiID, []byte(EmischdStructBytes))
	if emiScherr != nil {
		fmt.Println("Could not save EMI schdule data to ledger", emiScherr)
		return "", fmt.Errorf("Could not save EMI schdule data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "EMIScheduleIndex"
	EMIScheduleKey, err := stub.CreateCompositeKey(indexName, []string{EMIScheduleStruct.LoanID, EMIScheduleStruct.EmiID})
	if err != nil {
		fmt.Println("Could not index composite key for EMI Id ", err)
		return "", fmt.Errorf("Could not index composite key for EMI Id")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(EMIScheduleKey, value)

	indexName2 := "EMIhashIndex"
	EMIhashIndexKey, err := stub.CreateCompositeKey(indexName2, []string{EMIScheduleStruct.EMIHash, EMIScheduleStruct.EmiID})
	if err != nil {
		fmt.Println("Could not index composite key for EMIid~EMIHash", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value2 := []byte{0x00}
	stub.PutState(EMIhashIndexKey, value2)

	var emiSchdEvent = "{eventType: 'CreateEMISchedule', description:" + EMIScheduleStruct.LoanID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(emiSchdEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved EMISchedule")
	return args[0], nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(EMISchedule))
	{
		fmt.Printf("Error starting EMISchedule chaincode: %s", err)
	}
}
