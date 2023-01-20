package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type TransactionDetails struct {
	TransactionId  string `json:"TransactionId"`
	PayerId        string `json:"PayerId"`
	PayerRole      string `json:"PayerRole"`
	PayeeId        string `json:"PayeeId"`
	PayeeRole      string `json:"PayeeRole"`
	AmountPaid     string `json:"AmountPaid"`
	PaymentDate    string `json:"PaymentDate"`
	PaymentStatus  string `json:"PaymentStatus"`
	ApprovalStatus string `json:"ApprovalStatus"`
	Currency       string `json:"Currency"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *TransactionDetails) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *TransactionDetails) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SaveTransactionDetails" {

		result, err = SaveTransactionDetails(stub, args)
	}
	if fn == "GetTransactionDetailsByPayeeId" {

		result, err = GetTransactionDetailsByPayeeId(stub, args)
	}
	if fn == "GetTransactionDetailsByPayerId" {

		result, err = GetTransactionDetailsByPayerId(stub, args)
	}
	if fn == "GetTransactionDetailsByTransactionId" {

		result, err = GetTransactionDetailsByTransactionId(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetTransactionDetailsByTransactionId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTransactionDetailsByTransactionId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var transactionId = args[0]
	value, err := stub.GetState(transactionId)
	if err != nil {
		fmt.Println("Couldn't get details for this transaction id "+transactionId+" from ledger", err)
		return "", fmt.Errorf("Missing Transaction id")
	}
	return string(value), nil
}
func SaveTransactionDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveTransactionDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var transactionId = args[0]
	fmt.Println("the Transaction ID is" + transactionId)
	var payerId = args[1]
	fmt.Println("the Payer Id is" + payerId)
	var payerRole = args[2]
	fmt.Println("the Payer Role is" + payerRole)
	var payeeId = args[3]
	fmt.Println("the Payee Id is" + payeeId)
	var payeeRole = args[4]
	fmt.Println("the Payee Role is" + payeeRole)
	var amountPaid = args[5]
	fmt.Println("the Amount Paid is" + amountPaid)
	var paymentDate = args[6]
	fmt.Println("the Payment Date is" + paymentDate)
	var paymentStatus = args[7]
	fmt.Println("the Payment Status is" + paymentStatus)
	var approvalStatus = args[8]
	fmt.Println("the Approval Status is" + approvalStatus)
	var currency = args[9]
	fmt.Println("the currency is" + currency)

	//assigning to struct the variables
	TransactionDetailsStruct := TransactionDetails{
		TransactionId:  transactionId,
		PayerId:        payerId,
		PayerRole:      payerRole,
		PayeeId:        payeeId,
		PayeeRole:      payeeRole,
		AmountPaid:     amountPaid,
		PaymentDate:    paymentDate,
		PaymentStatus:  paymentStatus,
		ApprovalStatus: approvalStatus,
		Currency:       currency,
	}

	fmt.Println("the struct values are Transaction ID: " + TransactionDetailsStruct.TransactionId + "PayerId:" + TransactionDetailsStruct.PayerId + "PayeeId:" + TransactionDetailsStruct.PayeeId + "Payment Date: " + TransactionDetailsStruct.PaymentDate)

	transactionDetailsStructBytes, err := json.Marshal(TransactionDetailsStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	transactionerr := stub.PutState(transactionId, []byte(transactionDetailsStructBytes))
	if transactionerr != nil {
		fmt.Println("Could not save TransactionDetails data to ledger", transactionerr)
		return "", fmt.Errorf("Could not save TransactionDetails data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName1 := "transactionPayeeIndex"
	transactionPayeeIndexKey, err := stub.CreateCompositeKey(indexName1, []string{TransactionDetailsStruct.PayeeId, TransactionDetailsStruct.TransactionId})
	if err != nil {
		fmt.Println("Could not index composite key for transaction payee", err)
		return "", fmt.Errorf("Could not index composite key for transaction payee")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value1 := []byte{0x00}
	transactionerr = stub.PutState(transactionPayeeIndexKey, value1)
	if transactionerr != nil {
		fmt.Println("Could not save composite key to ledger", transactionerr)
		return "", fmt.Errorf("Could not save data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName2 := "transactionPayerIndex"
	transactionPayerIndexKey, err := stub.CreateCompositeKey(indexName2, []string{TransactionDetailsStruct.PayerId, TransactionDetailsStruct.TransactionId})
	if err != nil {
		fmt.Println("Could not index composite key for transaction payer", err)
		return "", fmt.Errorf("Could not index composite key for transaction payer")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value2 := []byte{0x00}
	transactionErr1 := stub.PutState(transactionPayerIndexKey, value2)
	if transactionErr1 != nil {
		fmt.Println("Could not save composite key to ledger", transactionerr)
		return "", fmt.Errorf("Could not save data to ledger")
	}

	var transactionEvent = "{eventType: 'SaveTransactionDetails', description:" + TransactionDetailsStruct.TransactionId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(transactionEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved Transaction Details")
	return args[0], nil

}

func GetTransactionDetailsByPayeeId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTransactionDetailsByPayeeId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var payeeId = args[0]

	transactionPayeeIndexItr, err := stub.GetStateByPartialCompositeKey("transactionPayeeIndex", []string{payeeId})
	if err != nil {
		return "", fmt.Errorf("Could not get Transaction for Payee")
	}
	defer transactionPayeeIndexItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{\"TransactionDetailsByPayeeId\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for transactionPayeeIndexItr.HasNext() {
		queryResponse, err := transactionPayeeIndexItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Transaction for Payee")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Cound not get composite key parts")
		}
		returnedPayeeId := compositeKeyParts[0]
		returnedTransactionId := compositeKeyParts[1]

		fmt.Printf(" - Value from index:%s TransactionId:%s PayeeId:%s\n", objectType, returnedTransactionId, returnedPayeeId)

		value, err := stub.GetState(returnedTransactionId)
		if err != nil {
			fmt.Println("Couldn't get Payee for TransactionId"+returnedTransactionId+"from ledger", err)
			return "", fmt.Errorf("Missing PayeeId")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	buffer.WriteString("}")

	fmt.Printf("- GetTrancheInvestorMapByTrancheId queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

func GetTransactionDetailsByPayerId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTransactionDetailsByPayerId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var payerId = args[0]

	transactionPayerIndexItr, err := stub.GetStateByPartialCompositeKey("transactionPayerIndex", []string{payerId})
	if err != nil {
		return "", fmt.Errorf("Could not get Transaction for Payer")
	}
	defer transactionPayerIndexItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"TransactionDetailsByPayerId\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for transactionPayerIndexItr.HasNext() {
		queryResponse, err := transactionPayerIndexItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Transaction for Payer")
		}

		// get transacitonId and payerId from transactionId~payerId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not split composite key transactionId~payerId")
		}
		returnedpayerId := compositeKeyParts[0]
		returnedTransactionId := compositeKeyParts[1]

		fmt.Printf(" - Found from index:%s payerID:%s TransactionID:%s\n", objectType, returnedpayerId, returnedTransactionId)

		value, err := stub.GetState(returnedTransactionId)
		if err != nil {
			fmt.Println("Could not get Composite Key Parts from ledger", err)
			return "", fmt.Errorf("Missing PayerId")
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	buffer.WriteString("}")
	fmt.Printf("- GetTransactionDetailsByPayerId queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(TransactionDetails))
	{
		fmt.Printf("Error starting TransactionDetails chaincode: %s", err)
	}
}
