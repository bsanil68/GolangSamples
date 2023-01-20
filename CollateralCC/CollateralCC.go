package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type LoanAsset struct {
	AssetID                       string `json:"AssetID"`
	LoanID                        string `json:"LoanID"`
	ClassOfVehicle                string `json:"ClassOfVehicle"`
	MakersName                    string `json:"MakersName"`
	TypeOfBody                    string `json:"TypeOfBody"`
	Horsepower                    string `json:"Horsepower"`
	ChassisNumber                 string `json:"ChassisNumber"`
	NumberOfCylinders             string `json:"NumberOfCylinders"`
	YearOfManufacturer            string `json:"YearOfManufacturer"`
	EngineNumber                  string `json:"EngineNumber"`
	Colour                        string `json:"Colour"`
	RegistrationNumber            string `json:"RegistrationNumber"`
	DateOfRTOGrantingRegistration string `json:"DateOfRTOGrantingRegistration"`
	Status                        string `json:"Status"`
	CollateralHash                string `json:"CollateralHash"`
}

func (t *LoanAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *LoanAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "CreateLoanAsset" {

		result, err = CreateLoanAsset(stub, args)
	}

	if fn == "GetLoansByAssetId" {

		result, err = GetLoansByAssetId(stub, args)
	}

	if fn == "GetLoansByCollateralHash" {

		result, err = GetLoansByCollateralHash(stub, args)
	}

	if fn == "GetLoansByLoanId" {

		result, err = GetLoansByLoanId(stub, args)
	}

	if fn == "GetAllLoans" {

		result, err = GetAllLoans(stub)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetLoansByAssetId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoansByAssetId")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing Asset ID")
	}
	// buffer is a JSON array containing QueryResults
	var assetID = args[0]
	value, err := stub.GetState(assetID)
	if err != nil {
		fmt.Println("Could not get Assets with id "+assetID+" from ledger", err)
		return "", errors.New("Missing AssetId")
	}
	//buffer.WriteString("}")
	return string(value), nil
}

func GetLoansByCollateralHash(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoansByCollateralHash")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing CollateralHash")
	}

	var collateral_Hash = args[0]

	getLoansCollateralHashItr, err := stub.GetStateByPartialCompositeKey("hashAssetIndex", []string{collateral_Hash})
	if err != nil {
		return "", fmt.Errorf("Could not get Loan Details for this collateral_Hash")
	}
	defer getLoansCollateralHashItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"AssetsForCollateralHash\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getLoansCollateralHashItr.HasNext() {
		queryResponse, err := getLoansCollateralHashItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Loan Data")
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
		returnedCollateralHash := compositeKeyParts[0]
		returnedAssetID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s collateralhash:%s assetid:%s\n", objectType, returnedCollateralHash, returnedAssetID)

		value, err := stub.GetState(returnedAssetID)
		if err != nil {
			fmt.Println("Couldn't get Loan Details for id "+returnedAssetID+" from ledger", err)
			return "", errors.New("Missing AssetID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetLoansByLoanId(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPaymentWaterfallDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing loanID")
	}

	var loanID = args[0]

	getLoansItr, err := stub.GetStateByPartialCompositeKey("loanAssetIndex", []string{loanID})
	if err != nil {
		return "", fmt.Errorf("Could not get Loan Details for this loan ID")
	}
	defer getLoansItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"AssetsForLoan\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getLoansItr.HasNext() {
		queryResponse, err := getLoansItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Loan Data")
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
		returnedLoanID := compositeKeyParts[0]
		returnedAssetID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s loanid:%s assetid:%s\n", objectType, returnedLoanID, returnedAssetID)

		value, err := stub.GetState(returnedAssetID)
		if err != nil {
			fmt.Println("Couldn't get Loan Details for id "+returnedAssetID+" from ledger", err)
			return "", errors.New("Missing AssetID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetAllLoans(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllLoans")

	getAllLoansItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Loans")
	}
	defer getAllLoansItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getAllLoansItr.HasNext() {
		queryResponseValue, err := getAllLoansItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Loan Data")
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

func CreateLoanAsset(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CreateLoanAsset")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	// json args input musts match  the case and  spelling exactly
	// get all the  arguments
	var assetID = args[0]
	fmt.Println("the asset id is" + assetID)
	var loanID = args[1]
	fmt.Println("the loan id is" + loanID)
	var classVehicle = args[2]
	fmt.Println("the classVehicle is" + classVehicle)
	var makersName = args[3]
	fmt.Println("the makersName is" + makersName)
	var typeOfBody = args[4]
	fmt.Println("the typeOfBody is" + typeOfBody)
	var horsePow = args[5]
	fmt.Println("the horsePow is" + horsePow)
	var chassisNum = args[6]
	fmt.Println("the chassisNum is" + chassisNum)
	var numOfCylinders = args[7]
	fmt.Println("the numOfCylinders is" + numOfCylinders)
	var yearOfManuf = args[8]
	fmt.Println("the yearOfManuf is" + yearOfManuf)
	var engineNum = args[9]
	fmt.Println("the engineNum is" + engineNum)
	var colour = args[10]
	fmt.Println("the Colour is" + colour)
	var registrationNum = args[11]
	fmt.Println("the registrationNum is" + registrationNum)
	var dateRTOGrantReg = args[12]
	fmt.Println("the dateRTOGrantReg is" + dateRTOGrantReg)
	var collateralHash = args[13]
	fmt.Println("the collateralHash is" + collateralHash)

	//assigning to struct the variables
	LoanAssetStruct := LoanAsset{
		AssetID:                       assetID,
		LoanID:                        loanID,
		ClassOfVehicle:                classVehicle,
		MakersName:                    makersName,
		TypeOfBody:                    typeOfBody,
		Horsepower:                    horsePow,
		ChassisNumber:                 chassisNum,
		NumberOfCylinders:             numOfCylinders,
		YearOfManufacturer:            yearOfManuf,
		EngineNumber:                  engineNum,
		Colour:                        colour,
		RegistrationNumber:            registrationNum,
		DateOfRTOGrantingRegistration: dateRTOGrantReg,
		Status:                        "active",
		CollateralHash:                collateralHash,
	}

	//  show that  the struct  has  values false

	fmt.Println("the struct values are assetID  " + LoanAssetStruct.AssetID + " loanID is " + LoanAssetStruct.LoanID + "and ClassOfVehicle is " + LoanAssetStruct.ClassOfVehicle)

	loanAssetStructBytes, err := json.Marshal(LoanAssetStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	loanAssetCaterr := stub.PutState(assetID, []byte(loanAssetStructBytes))
	if loanAssetCaterr != nil {
		fmt.Println("Could not save customer data to ledger", loanAssetCaterr)
		return "", fmt.Errorf("Could not save Loan asset data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "loanAssetIndex"
	loanAssetIndexKey, err := stub.CreateCompositeKey(indexName, []string{LoanAssetStruct.LoanID, LoanAssetStruct.AssetID})
	if err != nil {
		fmt.Println("Could not index composite key for collection", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(loanAssetIndexKey, value)

	indexName2 := "hashAssetIndex"
	hashAssetIndexKey, err := stub.CreateCompositeKey(indexName2, []string{LoanAssetStruct.CollateralHash, LoanAssetStruct.AssetID})
	if err != nil {
		fmt.Println("Could not index composite key for collection", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value2 := []byte{0x00}
	stub.PutState(hashAssetIndexKey, value2)

	var loanassetEvent = "{eventType: 'CreateLoanAsset', description: CollateralHash: " + LoanAssetStruct.CollateralHash + "AssetID" + LoanAssetStruct.AssetID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(loanassetEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved LoanAsset")
	return args[0], nil

}

func main() {
	err := shim.Start(new(LoanAsset))
	{
		fmt.Printf("Error starting LoanAsset chaincode: %s", err)
	}
}
