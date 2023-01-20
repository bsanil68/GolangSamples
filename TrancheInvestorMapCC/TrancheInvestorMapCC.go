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

type TrancheInvestorMap struct {
	TrancheId                   string `json:"TrancheId"`
	InvestorId                  string `json:"InvestorId"`
	NumberOfUnitsPurchased      string `json:"NumberOfUnitsPurchased"`
	AmountPaidByInvestor        string `json:"AmountPaidByInvestor"`
	ApplicationDate             string `json:"ApplicationDate"`
	AllocationDate              string `json:"AllocationDate"`
	MonthlyInterestForInvestor  string `json:"MonthlyInterestForInvestor"`
	MonthlyPrincipalForInvestor string `json:"MonthlyPrincipalForInvestor"`
	NumberOfInstallments        string `json:"NumberOfInstallments"`
	Currency                    string `json:"Currency"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *TrancheInvestorMap) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *TrancheInvestorMap) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SaveTrancheInvestorMap" {

		result, err = SaveTrancheInvestorMap(stub, args)
	}
	if fn == "GetTrancheInvestorMapByTrancheId" {

		result, err = GetTrancheInvestorMapByTrancheId(stub, args)
	}
	if fn == "GetTotalNumOfUnitsPurchasedOfTranche" {

		result, err = GetTotalNumOfUnitsPurchasedOfTranche(stub, args)
	}
	if fn == "GetTotalMonthlyInterestToBePaidToTrancheById" {

		result, err = GetTotalMonthlyInterestToBePaidToTrancheById(stub, args)
	}
	if fn == "GetTotalMonthlyPrincipalToBePaidToTrancheById" {

		result, err = GetTotalMonthlyPrincipalToBePaidToTrancheById(stub, args)
	}
	if fn == "CalculateMonthlyInterestForInvestor" {

		result, err = CalculateMonthlyInterestForInvestor(stub, args)
	}
	if fn == "CalculateMonthlyPrincipalForInvestor" {

		result, err = CalculateMonthlyPrincipalForInvestor(stub, args)
	}
	if fn == "GetTrancheInvestorMapByInvestorId" {

		result, err = GetTrancheInvestorMapByInvestorId(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func SaveTrancheInvestorMap(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveTrancheInvestorMap")

	if len(args) < 7 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var trancheId = args[0]
	fmt.Println("the Tranche ID is" + trancheId)
	var investorId = args[1]
	fmt.Println("the Investor ID is" + investorId)
	var numOfUnitsPurchased = args[2]
	fmt.Println("the Number of Units Purchased is" + numOfUnitsPurchased)
	var amountPaidByInvestor = args[3]
	fmt.Println("the Amount Paid by Investor is" + amountPaidByInvestor)
	var applicationDate = args[4]
	fmt.Println("the application date is" + applicationDate)
	var allocationDate = args[5]
	fmt.Println("the allocation date is" + allocationDate)
	var monthlyInterestForInvestor = args[6]
	fmt.Println("the Monthly Interest for Investor is" + monthlyInterestForInvestor)
	var monthlyPrincipalForInvestor = args[7]
	fmt.Println("the Monthly Principal for Investor is" + monthlyPrincipalForInvestor)
	var numOfInstallments = args[8]
	fmt.Println("the Number of Installments is" + numOfInstallments)
	var currency = args[9]
	fmt.Println("the currency is" + currency)

	//assigning to struct the variables
	TrancheInvestorMapStruct := TrancheInvestorMap{
		TrancheId:                   trancheId,
		InvestorId:                  investorId,
		NumberOfUnitsPurchased:      numOfUnitsPurchased,
		AmountPaidByInvestor:        amountPaidByInvestor,
		ApplicationDate:             applicationDate,
		AllocationDate:              allocationDate,
		MonthlyInterestForInvestor:  monthlyInterestForInvestor,
		MonthlyPrincipalForInvestor: monthlyPrincipalForInvestor,
		NumberOfInstallments:        numOfInstallments,
		Currency:                    currency,
	}

	fmt.Println("the struct values are Tranche ID: " + TrancheInvestorMapStruct.TrancheId + " Investor ID: " + TrancheInvestorMapStruct.InvestorId)

	trancheInvestorMapStructBytes, err := json.Marshal(TrancheInvestorMapStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	trancheInvestorErr := stub.PutState(investorId, []byte(trancheInvestorMapStructBytes))
	if trancheInvestorErr != nil {
		fmt.Println("Could not save TrancheInvestorMap data to ledger", trancheInvestorErr)
		return "", fmt.Errorf("Could not save TrancheInvestorMap data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "trancheInvestorMapIndex"
	trancheInvestorMapIndexKey, err := stub.CreateCompositeKey(indexName, []string{TrancheInvestorMapStruct.TrancheId, TrancheInvestorMapStruct.InvestorId})
	if err != nil {
		fmt.Println("Could not index composite key for tranche investor map", err)
		return "", fmt.Errorf("Could not index composite key for tranche investor map")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}

	trancheInvestorerr := stub.PutState(trancheInvestorMapIndexKey, value)
	if trancheInvestorerr != nil {
		fmt.Println("Cound not save data to ledger", trancheInvestorerr)
		return "", fmt.Errorf("Could not save data to ledger")
	}

	var trancheInvestorEvent = "{eventType: 'SaveTrancheInvestorMap', description:" + TrancheInvestorMapStruct.TrancheId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(trancheInvestorEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved Payment Tranche Investor Details")
	return args[0], nil
}

func GetTrancheInvestorMapByInvestorId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheInvestorMapByInvestorId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var investorID = args[0]
	value, err := stub.GetState(investorID)
	if err != nil {
		fmt.Println("Couldn't get details for this transaction id "+investorID+" from ledger", err)
		return "", fmt.Errorf("Missing Investor id")
	}
	return string(value), nil
}

func GetTrancheInvestorMapByTrancheId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheInvestorMapByTrancheId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var trancheId = args[0]

	trancheInvestorItr, err := stub.GetStateByPartialCompositeKey("trancheInvestorMapIndex", []string{trancheId})
	if err != nil {
		return "", fmt.Errorf("Could not get Investor Ids for Tranche")
	}
	defer trancheInvestorItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"TrancheInvestorMap\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for trancheInvestorItr.HasNext() {
		queryResponse, err := trancheInvestorItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Investor Ids for Tranche")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Cound not get composite key parts")
		}

		returnedTrancheID := compositeKeyParts[0]
		returnedInvestorID := compositeKeyParts[1]

		fmt.Printf(" - Value from index:%s InvestorId:%s TrancheId:%s\n", objectType, returnedInvestorID, returnedTrancheID)

		value, err := stub.GetState(returnedInvestorID)
		if err != nil {
			fmt.Println("Couldn't get Investor for TrancheID"+returnedInvestorID+"from ledger", err)
			return "", errors.New("Missing InvestorID")
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

func GetTotalNumOfUnitsPurchasedOfTranche(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTotalNumOfUnitsPurchasedOfTranche")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var trancheId = args[0]

	trancheInvestorItr, err := stub.GetStateByPartialCompositeKey("trancheInvestorMapIndex", []string{trancheId})
	if err != nil {
		return "", fmt.Errorf("Could not get Investor Ids for Tranche")
	}

	defer trancheInvestorItr.Close()

	var totalNumOfUnitsPurchasedOfTranche int
	totalNumOfUnitsPurchasedOfTranche = 0

	for trancheInvestorItr.HasNext() {
		queryResponse, err := trancheInvestorItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Investor Ids for Tranche")
		}

		// get the trancheId and investorId from trancheId~investorId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not split composite key trancheId~investorId")
		}
		returnedtrancheID := compositeKeyParts[0]
		returnedinvestorID := compositeKeyParts[1]

		fmt.Printf(" - Found from index:%s trancheID:%s investorID:%s\n", objectType, returnedtrancheID, returnedinvestorID)
		value, err := stub.GetState(returnedinvestorID)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}

		var buffer bytes.Buffer
		buffer.WriteString("{")
		buffer.WriteString(string(value))
		buffer.WriteString("}")

		var value1 = buffer.String()
		var trancheInv TrancheInvestorMap
		json.Unmarshal([]byte(value1), &trancheInv)
		numOfUnitsInvestor := trancheInv.NumberOfUnitsPurchased
		fmt.Printf(" Unmarshalled NumberOfUnitsPurchased:%s", numOfUnitsInvestor)
		var noOfUnitsInvestor int
		noOfUnitsInvestor, err = strconv.Atoi(numOfUnitsInvestor)
		totalNumOfUnitsPurchasedOfTranche += noOfUnitsInvestor
	}

	fmt.Printf("- GetTrancheInvestorMapByTrancheId totalNumOfUnitsPurchasedOfTranche:%s\n", string(totalNumOfUnitsPurchasedOfTranche))

	return string(totalNumOfUnitsPurchasedOfTranche), nil
}

func GetTotalMonthlyInterestToBePaidToTrancheById(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTotalMonthlyInterestToBePaidToTrancheById")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var trancheId = args[0]

	trancheInvestorItr, err := stub.GetStateByPartialCompositeKey("trancheInvestorMapIndex", []string{trancheId})
	if err != nil {
		return "", fmt.Errorf("Could not get Investor Ids for Tranche")
	}
	defer trancheInvestorItr.Close()

	var totalMonthlyInterestToBePaidToTranche int
	totalMonthlyInterestToBePaidToTranche = 0

	for trancheInvestorItr.HasNext() {
		queryResponse, err := trancheInvestorItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Investor Ids for Tranche")
		}

		// get the trancheId and investorId from trancheId~investorId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not split composite key trancheId~investorId")
		}

		returnedtrancheID := compositeKeyParts[0]
		returnedinvestorID := compositeKeyParts[1]

		fmt.Printf(" - Found from index:%s trancheID:%s investorID:%s\n", objectType, returnedtrancheID, returnedinvestorID)
		value, err := stub.GetState(returnedinvestorID)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		var buffer bytes.Buffer
		buffer.WriteString("{")
		buffer.WriteString(string(value))
		buffer.WriteString("}")

		var trancheInv TrancheInvestorMap
		var value1 = buffer.String()
		json.Unmarshal([]byte(value1), &trancheInv)

		monthlyInterestForInvestor := trancheInv.MonthlyInterestForInvestor
		fmt.Printf(" Unmarshalled MonthlyInterestForInvestor:%s", monthlyInterestForInvestor)
		var mnthlyInterest int
		mnthlyInterest, err = strconv.Atoi(monthlyInterestForInvestor)
		totalMonthlyInterestToBePaidToTranche += mnthlyInterest
	}

	fmt.Printf("- GetTotalMonthlyInterestToBePaidToTrancheById - totalMonthlyInterestToBePaidToTranche:%s\n", string(totalMonthlyInterestToBePaidToTranche))
	return string(totalMonthlyInterestToBePaidToTranche), nil
}

func GetTotalMonthlyPrincipalToBePaidToTrancheById(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTotalMonthlyPrincipalToBePaidToTrancheById")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var trancheId = args[0]

	trancheInvestorItr, err := stub.GetStateByPartialCompositeKey("trancheInvestorMapIndex", []string{trancheId})
	if err != nil {
		return "", fmt.Errorf("Could not get Investor Ids for Tranche")
	}
	defer trancheInvestorItr.Close()

	var totalMonthlyPrincipalToBePaidToTranche int
	totalMonthlyPrincipalToBePaidToTranche = 0

	for trancheInvestorItr.HasNext() {
		queryResponse, err := trancheInvestorItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get Investor Ids for Tranche")
		}
		// get the trancheId and investorId from trancheId~investorId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not split composite key trancheId~investorId")
		}

		returnedtrancheID := compositeKeyParts[0]
		returnedinvestorID := compositeKeyParts[1]

		fmt.Printf(" - Found from index:%s trancheID:%s investorID:%s\n", objectType, returnedtrancheID, returnedinvestorID)
		value, err := stub.GetState(returnedinvestorID)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		var buffer bytes.Buffer
		buffer.WriteString("{")
		buffer.WriteString(string(value))
		buffer.WriteString("}")

		var trancheInv TrancheInvestorMap
		var value1 = buffer.String()
		json.Unmarshal([]byte(value1), &trancheInv)

		monthlyPrincipalForInvestor := trancheInv.MonthlyPrincipalForInvestor
		fmt.Printf(" Unmarshalled MonthlyPrincipalForInvestor:%s", monthlyPrincipalForInvestor)
		var mnthlyPrincipalForInvestor int
		mnthlyPrincipalForInvestor, err = strconv.Atoi(monthlyPrincipalForInvestor)
		totalMonthlyPrincipalToBePaidToTranche += mnthlyPrincipalForInvestor
	}

	fmt.Printf("- GetTotalMonthlyPrincipalToBePaidToTrancheById -totalMonthlyPrincipalToBePaidToTranche :%s\n", string(totalMonthlyPrincipalToBePaidToTranche))
	return string(totalMonthlyPrincipalToBePaidToTranche), nil
}

func CalculateMonthlyInterestForInvestor(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CalculateMonthlyInterestForInvestor")

	if len(args) < 3 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Invalid number of arguments")
	}

	str := args[0]
	numOfUnitsPurchased, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of units purchased")
	}
	str1 := args[1]
	faceValue, err := strconv.ParseInt(str1, 10, 64)
	if err != nil {
		fmt.Println("Invalid face value")
	}
	str2 := args[2]
	couponRate, err := strconv.ParseInt(str2, 10, 64)
	if err != nil {
		fmt.Println("Invalid coupon Rate")
	}

	var annualInterest = (faceValue * numOfUnitsPurchased * couponRate) / 100

	var monthlyInterest = annualInterest / 12

	return string(monthlyInterest), nil
}

func CalculateMonthlyPrincipalForInvestor(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CalculateMonthlyPrincipalForInvestor")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Invalid number of arguments")
	}

	str := args[0]
	numOfUnitsPurchased, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of units purchased")
	}
	str1 := args[1]
	faceValue, err := strconv.ParseInt(str1, 10, 64)
	if err != nil {
		fmt.Println("Invalid face value")
	}

	var totalPrincipal = faceValue * numOfUnitsPurchased

	var monthlyPrincipal = totalPrincipal / 12

	return string(monthlyPrincipal), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(TrancheInvestorMap))
	{
		fmt.Printf("Error starting TrancheInvestorMap chaincode: %s", err)
	}
}
