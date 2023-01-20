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

type TrancheCC struct {
	PoolID                   string `json:"PoolID"`
	TrancheID                string `json:"TrancheID"`
	ClassName                string `json:"ClassName"`
	PercentOfCollateralValue string `json:"PercentOfCollateralValue"`
	NominalAmount            string `json:"NominalAmount"`
	IssuePricePercent        string `json:"IssuePricePercent"`
	Size                     string `json:"Size"`
	InterestRateType         string `json:"InterestRateType"`
	CouponRate               string `json:"CouponRate"`
	InterestFrequency        string `json:"InterestFrequency"`
	FaceValue                string `json:"FaceValue"`
	MinUnitsOfSubscription   string `json:"MinUnitsOfSubscription"`
	TotalUnitsAvailable      string `json:"TotalUnitsAvailable"`
	NotePeriod               string `json:"NotePeriod"`
	NoOfUnitsLeft		 string `json:"NoOfUnitsLeft"`
	Currency                 string `json:"Currency"`
}

func (t *TrancheCC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *TrancheCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetTrancheByPoolID" {

		result, err = GetTrancheByPoolID(stub, args)
	}
	if fn == "GetTrancheByTrancheID" {

		result, err = GetTrancheByTrancheID(stub, args)
	}
	if fn == "GetAllTranches" {

		result, err = GetAllTranches(stub)
	}
	if fn == "SaveTranche" {

		result, err = SaveTranche(stub, args)
	}
	if fn == "CalculateAmountPayableForUnitsPurchased" {

		result, err = CalculateAmountPayableForUnitsPurchased(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetTrancheByTrancheID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheByTrancheID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing trancheId")
	}

	var tranche_Id = args[0]
	value, err := stub.GetState(tranche_Id)
	if err != nil {
		fmt.Println("Couldn't get trancheId  with id "+tranche_Id+" from ledger", err)
		return "", errors.New("Missing trancheId")
	}

	return string(value), nil
}

func GetTrancheByPoolID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheByPoolID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing TrancheID")
	}

	var poolId = args[0]

	getTrancheByPoolIDItr, err := stub.GetStateByPartialCompositeKey("trancheIdIndex", []string{poolId})
	if err != nil {
		return "", fmt.Errorf("Could not get tranches for this pool ID")
	}
	defer getTrancheByPoolIDItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"TrancheForPoolID\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getTrancheByPoolIDItr.HasNext() {
		queryResponse, err := getTrancheByPoolIDItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next tranche Data")
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
		returnedtrancheID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s trancheid:%s poolid:%s\n", objectType, returnedtrancheID, returnedPoolID)

		value, err := stub.GetState(returnedtrancheID)
		if err != nil {
			fmt.Println("Couldn't get tranche for id "+returnedtrancheID+" from ledger", err)
			return "", errors.New("Missing trancheID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetAllTranches(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllTranches")

	getAllTranchesItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all tranches")
	}
	defer getAllTranchesItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"GetAllTranches\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getAllTranchesItr.HasNext() {
		queryResponseValue, err := getAllTranchesItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Pool Data")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponseValue.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func SaveTranche(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveTranche")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var poolID = args[0]
	fmt.Println("the PoolID is" + poolID)
	var trancheID = args[1]
	fmt.Println("the TrancheID is" + trancheID)
	var className = args[2]
	fmt.Println("the class Name is" + className)
	var percentOfCollateralValue = args[3]
	fmt.Println("the PercentOfCollateralValue is" + percentOfCollateralValue)
	var nominalAmount = args[4]
	fmt.Println("the nominal Amount is" + nominalAmount)
	var issuePricePercent = args[5]
	fmt.Println("the IssuePricePercent is" + issuePricePercent)
	var Size = args[6]
	fmt.Println("the Size is" + Size)
	var interestRateType = args[7]
	fmt.Println("the InterestRateType is" + interestRateType)
	var couponRate = args[8]
	fmt.Println("the CouponRate is" + couponRate)
	var interestFrequency = args[9]
	fmt.Println("the InterestFrequency is" + interestFrequency)
	var faceValue = args[10]
	fmt.Println("the FaceValue is" + faceValue)
	var minUnitsOfSubscription = args[11]
	fmt.Println("the MinUnitsOfSubscription is" + minUnitsOfSubscription)
	var totalUnitsAvailable = args[12]
	fmt.Println("the TotalUnitsAvailable is" + totalUnitsAvailable)
	var notePeriod = args[13]
	fmt.Println("the NotePeriod is" + notePeriod)
	var noOfUnitsLeft = args[14]
	fmt.Println("the noOfUnitsLeft is" + noOfUnitsLeft)
	var currency = args[15]
	fmt.Println("the NotePeriod is" + currency)

	TrancheCCStruct := TrancheCC{
		PoolID:                   poolID,
		TrancheID:                trancheID,
		ClassName:                className,
		PercentOfCollateralValue: percentOfCollateralValue,
		NominalAmount:            nominalAmount,
		IssuePricePercent:        issuePricePercent,
		Size:                     Size,
		InterestRateType:         interestRateType,
		CouponRate:               couponRate,
		InterestFrequency:        interestFrequency,
		FaceValue:                faceValue,
		MinUnitsOfSubscription:   minUnitsOfSubscription,
		TotalUnitsAvailable:      totalUnitsAvailable,
		NotePeriod:               notePeriod,
		NoOfUnitsLeft:		  noOfUnitsLeft,
		Currency:                 currency,
	}

	fmt.Println("the struct values are PoolID  " + TrancheCCStruct.PoolID + "TrancheID" + TrancheCCStruct.TrancheID + "ClassName" + TrancheCCStruct.ClassName + "PercentOfCollateralValue" + TrancheCCStruct.PercentOfCollateralValue + "NominalAmount" + TrancheCCStruct.NominalAmount + "IssuePricePercent" + TrancheCCStruct.IssuePricePercent + "Size" + TrancheCCStruct.Size + "InterestRateType" + TrancheCCStruct.InterestRateType + "CouponRate" + TrancheCCStruct.CouponRate + "InterestFrequency" +
		TrancheCCStruct.InterestFrequency + "FaceValue  " + TrancheCCStruct.FaceValue + "MinUnitsOfSubscription" + TrancheCCStruct.MinUnitsOfSubscription + "TotalUnitsAvailable  " + TrancheCCStruct.TotalUnitsAvailable + "NotePeriod" + TrancheCCStruct.NotePeriod)

	trancheStructBytes, err := json.Marshal(TrancheCCStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	trancheErr := stub.PutState(trancheID, []byte(trancheStructBytes))
	if trancheErr != nil {
		fmt.Println("Couldn't save trancheCharacterestic data to ledger", trancheErr)
		return "", fmt.Errorf("Couldn't save trancheCharacterestic data to ledger")
	}

	indexName := "trancheIdIndex"
	trancheIdKey, err := stub.CreateCompositeKey(indexName, []string{TrancheCCStruct.PoolID, TrancheCCStruct.TrancheID})
	if err != nil {
		fmt.Println("Could not index composite key for trancheId", err)
		return "", fmt.Errorf("Could not index composite key for TrancheID map")
	}

	value := []byte{0x00}
	stub.PutState(trancheIdKey, value)

	var trancheEvent = "{eventType: 'TrancheCC', description:" + TrancheCCStruct.PoolID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(trancheEvent))
	if err != nil {
		return "", fmt.Errorf("Couldn't set event hub")
	}
	fmt.Println("Successfully saved TrancheCC")
	return args[0], nil
}

func CalculateAmountPayableForUnitsPurchased(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CalculateAmountPayableForUnitsPurchased")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		//return "", errors.New("Missing userId")
	}

	str := args[0]
	numberofunitspur, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Invalid number of units purchased")
	}
	str1 := args[1]
	facevalue, err := strconv.ParseInt(str1, 10, 64)
	if err != nil {
		fmt.Println("Invalid face value")
	}
	var amountPayable = facevalue * numberofunitspur

	return string(amountPayable), nil
}

func main() {
	err := shim.Start(new(TrancheCC))
	{
		fmt.Printf("Error starting TrancheCC chaincode: %s", err)
	}
}
