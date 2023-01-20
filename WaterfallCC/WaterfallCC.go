package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type PaymentWaterfall struct {
	TransactionID                    string `json:"TransactionID"`
	PoolID                           string `json:"PoolID"`
	TaxPayment                       string `json:"TaxPayment"`
	PaymentOfReimbursableExpenses    string `json:"PaymentOfReimbursableExpenses"`
	PaymentOfServicingFeesToServicer string `json:"PaymentOfServicingFeesToServicer"`
	InterestPayableToClassA          string `json:"InterestPayableToClassA"`
	PrincipalPayableToClassA         string `json:"PrincipalPayableToClassA"`
	InterestPayableToClassB          string `json:"InterestPayableToClassB"`
	PrincipalPayableToClassB         string `json:"PrincipalPayableToClassB"`
	InterestPayableToClassC          string `json:"InterestPayableToClassC"`
	PrincipalPayableToClassC         string `json:"PrincipalPayableToClassC"`
	BalanceAsIncomeOfClassC          string `json:"BalanceAsIncomeOfClassC"`
	Month                            string `json:"Month"`
	TotalCollections                 string `json:"TotalCollections"`
	Currency                         string `json:"Currency"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *PaymentWaterfall) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *PaymentWaterfall) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SavePaymentWaterfallDetails" {

		result, err = SavePaymentWaterfallDetails(stub, args)
	}
	if fn == "GetPaymentWaterfallDetails" {

		result, err = GetPaymentWaterfallDetails(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetPaymentWaterfallDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPaymentWaterfallDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing TransactionID")
	}

	var poolID = args[0]
	var month = args[1]

	getPaymentWaterfallItr, err := stub.GetStateByPartialCompositeKey("poolMonthWaterfallIndex", []string{poolID, month})
	if err != nil {
		return "", fmt.Errorf("Could not get Payment Waterfall Details for this transaction ID")
	}
	defer getPaymentWaterfallItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"PaymentWaterfallDetailsForTransactionID\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getPaymentWaterfallItr.HasNext() {
		queryResponse, err := getPaymentWaterfallItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next PaymentWaterfallDetails Data")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		returnedPoolID := compositeKeyParts[0]
		returnedmonth := compositeKeyParts[1]
		returnedtransactionID := compositeKeyParts[2]

		fmt.Printf("- found a  from index:%s poolid:%s month:%s transactionID:%s\n", objectType, returnedPoolID, returnedmonth, returnedtransactionID)

		value, err := stub.GetState(returnedtransactionID)
		if err != nil {
			fmt.Println("Couldn't get PaymentWaterfallDetails for id "+returnedtransactionID+" from ledger", err)
			return "", errors.New("Missing transactionID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func SavePaymentWaterfallDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SavePaymentWaterfallDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	var transactionID = args[0]
	fmt.Println("the Transaction ID is" + transactionID)
	var poolID = args[1]
	fmt.Println("the Pool ID is" + poolID)
	var taxPayment = args[2]
	fmt.Println("the Tax Payment is" + taxPayment)
	var paymentOfReimbursableExpenses = args[3]
	fmt.Println("the Payment of Reimbursable Expenses is" + paymentOfReimbursableExpenses)
	var paymentOfServicingFeesToServicer = args[4]
	fmt.Println("the Payment of Servicing Fees to Servicer is" + paymentOfServicingFeesToServicer)
	var interestPayableToClassA = args[5]
	fmt.Println("the Interest Payable to Class A is" + interestPayableToClassA)
	var principalPayableToClassA = args[6]
	fmt.Println("the Interest Payable to Class A is" + principalPayableToClassA)
	var interestPayableToClassB = args[7]
	fmt.Println("the Interest Payable to Class B is" + interestPayableToClassB)
	var principalPayableToClassB = args[8]
	fmt.Println("the Interest Payable to Class B is" + principalPayableToClassB)
	var interestPayableToClassC = args[9]
	fmt.Println("the Interest Payable to Class C is" + interestPayableToClassC)
	var principalPayableToClassC = args[10]
	fmt.Println("the Interest Payable to Class C is" + principalPayableToClassC)
	var balanceAsIncomeOfClassC = args[11]
	fmt.Println("the Balance as Income of Class C is" + balanceAsIncomeOfClassC)
	var month = args[12]
	fmt.Println("the month of Collections is" + month)
	var totalCollections = args[13]
	fmt.Println("the Total Collections is" + totalCollections)
	var currency = args[14]
	fmt.Println("the currency is" + currency)

	//assigning to struct the variables
	PaymentWaterfallStruct := PaymentWaterfall{
		TransactionID:                    transactionID,
		PoolID:                           poolID,
		TaxPayment:                       taxPayment,
		PaymentOfReimbursableExpenses:    paymentOfReimbursableExpenses,
		PaymentOfServicingFeesToServicer: paymentOfServicingFeesToServicer,
		InterestPayableToClassA:          interestPayableToClassA,
		PrincipalPayableToClassA:         principalPayableToClassA,
		InterestPayableToClassB:          interestPayableToClassB,
		PrincipalPayableToClassB:         principalPayableToClassB,
		InterestPayableToClassC:          interestPayableToClassC,
		PrincipalPayableToClassC:         principalPayableToClassC,
		BalanceAsIncomeOfClassC:          balanceAsIncomeOfClassC,
		Month:                            month,
		TotalCollections:                 totalCollections,
		Currency:                         currency,
	}

	fmt.Println("the struct values are Transaction ID: " + PaymentWaterfallStruct.TransactionID + " Pool ID: " + PaymentWaterfallStruct.PoolID + " Tax Payment: " + PaymentWaterfallStruct.TaxPayment + " Payment Of Reimbursable Expenses: " + PaymentWaterfallStruct.PaymentOfReimbursableExpenses + " Payment Of Servicing Fees To Servicer: " + PaymentWaterfallStruct.PaymentOfServicingFeesToServicer + " Interest Payable To ClassA: " + PaymentWaterfallStruct.InterestPayableToClassA + " Principal Payable To ClassA: " + PaymentWaterfallStruct.PrincipalPayableToClassA + " Interest Payable To ClassB: " + PaymentWaterfallStruct.InterestPayableToClassB + " Principal Payable To ClassB: " + PaymentWaterfallStruct.PrincipalPayableToClassB + " InterestPayableToClassC: " + PaymentWaterfallStruct.InterestPayableToClassC + " Principal Payable To ClassC: " + PaymentWaterfallStruct.PrincipalPayableToClassC + " Balance As Income Of ClassC: " + PaymentWaterfallStruct.BalanceAsIncomeOfClassC + "Month: " + PaymentWaterfallStruct.Month + "Total Collections: " + PaymentWaterfallStruct.TotalCollections)

	paymentWaterfallStructBytes, err := json.Marshal(PaymentWaterfallStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	waterfallerr := stub.PutState(transactionID, []byte(paymentWaterfallStructBytes))
	if waterfallerr != nil {
		fmt.Println("Could not save PaymentWaterfall data to ledger", waterfallerr)
		return "", fmt.Errorf("Could not save PaymentWaterfall data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "poolMonthWaterfallIndex"
	poolMonthWaterfallIndexKey, err := stub.CreateCompositeKey(indexName, []string{PaymentWaterfallStruct.PoolID, PaymentWaterfallStruct.Month, PaymentWaterfallStruct.TransactionID})
	if err != nil {
		fmt.Println("Could not index composite key for collection", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(poolMonthWaterfallIndexKey, value)

	var waterfallEvent = "{eventType: 'SavePaymentWaterfallDetails', description:" + PaymentWaterfallStruct.PoolID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(waterfallEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved Payment Waterfall Details")
	return args[0], nil

}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(PaymentWaterfall))
	{
		fmt.Printf("Error starting PaymentWaterfall chaincode: %s", err)
	}
}
