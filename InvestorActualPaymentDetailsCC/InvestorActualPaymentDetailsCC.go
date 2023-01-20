package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type InvestorActualPaymentDetailsCC struct {
	InvestorID           string `json:"InvestorId"`
	InterestPaid         string `json:"InterestPaid"`
	MonthOfInterest      string `json:"MonthOfInterest"`
	PrincipalPaid        string `json:"PrincipalPaid"`
	MonthOfPrincipal     string `json:"MonthOfPrincipal"`
	OutstandingInterest  string `json:"OutstandingInterest"`
	OutstandingPrincipal string `json:"OutstandingPrincipal"`
	MonthPaidIn          string `json:"MonthPaidIn"`
	OverdueInterest      string `json:"OverdueInterest"`
	OverduePrincipal     string `json:"OverduePrincipal"`
	Currency             string `json:"Currency"`
}

func (t *InvestorActualPaymentDetailsCC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *InvestorActualPaymentDetailsCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error

	if fn == "GetInvestorPaymentDetails" {

		result, err = GetInvestorPaymentDetails(stub, args)
	}

	if fn == "GetInvestorPaymentDetailsByMonth" {

		result, err = GetInvestorPaymentDetailsByMonth(stub, args)
	}

	if fn == "SaveInvestorPaymentDetails" {

		result, err = SaveInvestorPaymentDetails(stub, args)
	}

	if fn == "CalculateOutstandingInterestOfInvestor" {

		result, err = CalculateOutstandingInterestOfInvestor(stub, args)
	}

	if fn == "CalculateOutstandingPrincipalOfInvestor" {

		result, err = CalculateOutstandingPrincipalOfInvestor(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetInvestorPaymentDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInvestorPaymentDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing InvestorId")
	}

	var Investor_ID = args[0]

	value, err := stub.GetState(Investor_ID)
	if err != nil {
		fmt.Println("Couldn't get InvestorId  with id "+Investor_ID+" from ledger", err)
		return "", errors.New("Missing Investor_Id")
	}

	return string(value), nil
}

func GetInvestorPaymentDetailsByMonth(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInvestorPaymentDetailsByMonth")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing Month")
	}

	var month = args[0]

	investorPaymentItr, err := stub.GetStateByPartialCompositeKey("InvestorActualPaymentIndex", []string{month})
	if err != nil {

		return "", fmt.Errorf("Could not get payment schedule for month")
	}

	defer investorPaymentItr.Close()

	//  the root  for  json so that  it cane unmarshalled  to pojo

	var buffer bytes.Buffer

	buffer.WriteString("{\"InvestorActualPaymentDetails\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	// Iterate through result set and check
	for investorPaymentItr.HasNext() {
		// Note that we don't get the value (2nd return variable), we'll just get the loan id from the composite key
		responseRange, err := investorPaymentItr.Next()
		if err != nil {
			return "", fmt.Errorf("no Loans assigned for pool")
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// get the investor id and month from InvestorActualPaymentIndex composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			fmt.Println("- split error   ")
			return "", fmt.Errorf("splitting error ")
		}

		returnedMonth := compositeKeyParts[0]
		returnedInvestorID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s loanid:%s poolid:%s\n", objectType, returnedMonth, returnedInvestorID)
		value, err := stub.GetState(returnedInvestorID)
		if err != nil {
			fmt.Println("Could not get assigner for id "+returnedMonth+" from ledger", err)
			return "", errors.New("Missing Month")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")
	return buffer.String(), nil
}

func SaveInvestorPaymentDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveInvestorPaymentDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var investorID = args[0]
	fmt.Println("the Investor Id is" + investorID)
	var interestPaid = args[1]
	fmt.Println("the Interest Paid is" + interestPaid)
	var monthOfInterest = args[2]
	fmt.Println("the Month Of Interest is" + monthOfInterest)
	var principalPaid = args[3]
	fmt.Println("the Principal Paid is" + principalPaid)
	var monthOfPrincipal = args[4]
	fmt.Println("the Month Of Principal is" + monthOfPrincipal)
	var outstandingInterest = args[5]
	fmt.Println("the Outstanding Interest is" + outstandingInterest)
	var outstandingPrincipal = args[6]
	fmt.Println("the OutstandingPrincipal is" + outstandingPrincipal)
	var monthPaidIn = args[7]
	fmt.Println("the Month Paid In is" + monthPaidIn)
	var overdueInterest = args[8]
	fmt.Println("the overdueInterest is" + overdueInterest)
	var overduePrincipal = args[9]
	fmt.Println("the overduePrincipal is" + overduePrincipal)
	var currency = args[10]
	fmt.Println("the currency is" + currency)

	InvestorActualPaymentDetailsCCStruct := InvestorActualPaymentDetailsCC{
		InvestorID:           investorID,
		InterestPaid:         interestPaid,
		MonthOfInterest:      monthOfInterest,
		PrincipalPaid:        principalPaid,
		MonthOfPrincipal:     monthOfPrincipal,
		OutstandingInterest:  outstandingInterest,
		OutstandingPrincipal: outstandingPrincipal,
		MonthPaidIn:          monthPaidIn,
		OverdueInterest:      overdueInterest,
		OverduePrincipal:     overduePrincipal,
		Currency:             currency,
	}

	fmt.Println("the struct values are InvestorId  " + InvestorActualPaymentDetailsCCStruct.InvestorID + "InterestPaid" + InvestorActualPaymentDetailsCCStruct.InterestPaid + "MonthOfInterest" + InvestorActualPaymentDetailsCCStruct.MonthOfInterest + "PrincipalPaid" + InvestorActualPaymentDetailsCCStruct.PrincipalPaid + "MonthOfPrincipal" + InvestorActualPaymentDetailsCCStruct.MonthOfPrincipal + "OutstandingInterest" + InvestorActualPaymentDetailsCCStruct.OutstandingInterest + "OutstandingPrincipal" + InvestorActualPaymentDetailsCCStruct.OutstandingPrincipal + "MonthPaidIn" + InvestorActualPaymentDetailsCCStruct.MonthPaidIn + "OverdueInterest" + InvestorActualPaymentDetailsCCStruct.OverdueInterest + "OverduePrincipal" + InvestorActualPaymentDetailsCCStruct.OverduePrincipal)

	ActualpaymentStructBytes, err := json.Marshal(InvestorActualPaymentDetailsCCStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}

	ActualpaymentStructerr := stub.PutState(investorID, []byte(ActualpaymentStructBytes))
	if ActualpaymentStructerr != nil {
		fmt.Println("Could not save Actualpayment data to ledger", ActualpaymentStructerr)
		return "", fmt.Errorf("Could not save Actualpayment data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "InvestorActualPaymentIndex"
	investorActualPaymentKey, err := stub.CreateCompositeKey(indexName, []string{InvestorActualPaymentDetailsCCStruct.MonthOfInterest, InvestorActualPaymentDetailsCCStruct.InvestorID})
	if err != nil {
		fmt.Println("Could not index composite key for pool loan", err)
		return "", fmt.Errorf("Could not index composite key for pool loan")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(investorActualPaymentKey, value)

	var ActualpaymentEvent = "{eventType: 'InvestorActualPaymentDetailsCCStruct', description:" + InvestorActualPaymentDetailsCCStruct.InvestorID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(ActualpaymentEvent))
	if err != nil {
		fmt.Println("Successfully saved Investor Actual Payment Details")
		return "", fmt.Errorf("Could not set event hub")
	}
	return args[0], nil

}

func CalculateOutstandingInterestOfInvestor(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CalculateOutstandingInterestOfInvestor")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing arguments")
	}

	str := args[0]
	collectionAmount, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Invalid Collection amount")
	}
	str1 := args[1]
	numberOfUnitsPur, err := strconv.ParseInt(str1, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of units purchased")
	}
	str2 := args[2]
	faceValue, err := strconv.ParseInt(str2, 10, 64)
	if err != nil {
		fmt.Println("Invalid face value")
	}
	str3 := args[3]
	couponRate, err := strconv.ParseInt(str3, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of coupon rate")
	}
	str4 := args[4]
	notePeriod, err := strconv.ParseInt(str4, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of note period")
	}

	var totalInterestForInvestor = (numberOfUnitsPur * faceValue) * (couponRate / 100)
	var totalPrincipalForInvestor = numberOfUnitsPur * faceValue
	var monthlyInterestForInvestor = totalInterestForInvestor / 12
	var monthlyPrincipalForInvestor = totalPrincipalForInvestor / 12
	var outstandingInterestOfInvestor = (totalInterestForInvestor * notePeriod) - monthlyInterestForInvestor

	if collectionAmount > monthlyPrincipalForInvestor {
		outstandingInterestOfInvestor -= monthlyInterestForInvestor
	} else {
		outstandingInterestOfInvestor -= collectionAmount
	}
	return string(outstandingInterestOfInvestor), nil
}

func CalculateOutstandingPrincipalOfInvestor(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CalculateOutstandingPrincipalOfInvestor")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing arguments")
	}

	str5 := args[0]
	collectionAmount, err := strconv.ParseInt(str5, 10, 64)
	if err != nil {
		fmt.Println("Invalid Collection amount")
	}
	str6 := args[1]
	numberOfUnitsPur, err := strconv.ParseInt(str6, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of units purchased")
	}
	str7 := args[2]
	faceValue, err := strconv.ParseInt(str7, 10, 64)
	if err != nil {
		fmt.Println("Invalid face value")
	}
	str8 := args[3]
	outstandingPrincipalOfPreviousmonth, err := strconv.ParseInt(str8, 10, 64)
	if err != nil {
		fmt.Println("Invalid outstandingPrincipalOfPreviousmonth")
	}
	var totalPrincipalForInvestor = numberOfUnitsPur * faceValue
	var monthlyPrincipalForInvestor = totalPrincipalForInvestor / 12

	if collectionAmount > monthlyPrincipalForInvestor {
		outstandingPrincipalOfPreviousmonth -= monthlyPrincipalForInvestor
	} else {
		outstandingPrincipalOfPreviousmonth -= collectionAmount
	}
	return string(outstandingPrincipalOfPreviousmonth), nil
}

func main() {
	err := shim.Start(new(InvestorActualPaymentDetailsCC))
	{
		fmt.Printf("Error starting InvestorActualPaymentDetailsCC chaincode: %s", err)
	}
}
