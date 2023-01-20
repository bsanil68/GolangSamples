package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type InvestorReport struct {
	IRID           string `json:"irId"`
	IRData         string `json:"irData"`
	IRHash         string `json:"irHash"`
	IRMonth        string `json:"irMonth"`
	IRYear         string `json:"irYear"`
	IRUpdatedBy    string `json:"irUpdatedBy"`
	IRUpdationDate string `json:"irUpdationDate"`
	IRDealID       string `json:"irDealId"`
}

func (t *InvestorReport) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *InvestorReport) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetInvestorReportByID" {

		result, err = GetInvestorReportByID(stub, args)
	}
	if fn == "GetInvestorReportByMonthAndYear" {

		result, err = GetInvestorReportByMonthAndYear(stub, args)
	}
	if fn == "SaveInvestorReport" {

		result, err = SaveInvestorReport(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func hash1(v []byte) string {
	h := sha256.New()
	h.Write(v)
	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}

func hash2_serialize_then_hash(v1 []string) string {
	bs, _ := json.Marshal(v1)
	return hash1(bs)
}

func GetInvestorReportByID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInvestorReportByID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing IRID")
	}

	var irId = args[0]
	value, err := stub.GetState(irId)
	if err != nil {
		fmt.Println("Couldn't get Investor Report with id "+irId+" from ledger", err)
		return "", errors.New("Missing irId")
	}

	return string(value), nil
}

func GetInvestorReportByMonthAndYear(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInvestorReportByMonthAndYear")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing month or year")
	}

	var irDealID = args[0]
	var month = args[1]
	var year = args[2]

	getIRByMonthAndYearItr, err := stub.GetStateByPartialCompositeKey("investorReportMonthIndex", []string{irDealID, month, year})
	if err != nil {
		return "", fmt.Errorf("Could not get Investor Report for this irDealID , Month and Year")
	}
	defer getIRByMonthAndYearItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	// buffer.WriteString("{\"DTOInvestorReportData\":")
	bArrayMemberAlreadyWritten := false
	for getIRByMonthAndYearItr.HasNext() {
		queryResponse, err := getIRByMonthAndYearItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next data")
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
		returnedIRDealID := compositeKeyParts[0]
		returnedIRMonth := compositeKeyParts[1]
		returnedIRYear := compositeKeyParts[2]
		returnedIRID := compositeKeyParts[3]

		fmt.Printf("- found a  from index:%s irDealID:%s irID:%s irMonth:%s irYear:%s\n", objectType, returnedIRDealID, returnedIRID, returnedIRMonth, returnedIRYear)

		value, err := stub.GetState(returnedIRID)
		if err != nil {
			fmt.Println("Couldn't get investor report for id "+returnedIRID+" from ledger", err)
			return "", errors.New("Missing irID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
		break
	}
	//buffer.WriteString("}")

	return buffer.String(), nil
}

func SaveInvestorReport(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveInvestorReport")

	if len(args) < 6 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var irID = args[0]
	fmt.Println("the IRID is" + irID)

	var jsonObj interface{}

	var buffer bytes.Buffer
	buffer.WriteString("\"")
	buffer.WriteString(args[1])
	buffer.WriteString("\"")

	fmt.Println("IRData", args[1])

	err := json.Unmarshal([]byte(buffer.String()), &jsonObj)
	if err != nil {
		fmt.Println("Can't deserialize", []byte(args[1]))
	}

	b, _ := json.Marshal(jsonObj)

	var irData = string(b)
	fmt.Println("the Investor Report Data is" + irData)
	var irHash = hash2_serialize_then_hash([]string{irData})
	fmt.Println("the Investor Report Hash is" + irHash)
	var irMonth = args[2]
	fmt.Println("the Month is" + irMonth)
	var irYear = args[3]
	fmt.Println("the Year is" + irYear)
	var irUpdatedBy = args[4]
	fmt.Println("the User ID is" + irUpdatedBy)
	var irUpdationDate = args[5]
	fmt.Println("the Updation Date is" + irUpdationDate)
	var irDealID = args[6]
	fmt.Println("the irDealID is" + irDealID)

	InvestorReportStruct := InvestorReport{
		IRID:           irID,
		IRData:         irData,
		IRHash:         irHash,
		IRMonth:        irMonth,
		IRYear:         irYear,
		IRUpdatedBy:    irUpdatedBy,
		IRUpdationDate: irUpdationDate,
		IRDealID:       irDealID,
	}

	fmt.Println("the struct values are IRID:" + InvestorReportStruct.IRID + " IRData:" + InvestorReportStruct.IRData +
		" IRMonth:" + InvestorReportStruct.IRMonth + " IRYear:" + InvestorReportStruct.IRYear)

	investorReportStructBytes, err := json.Marshal(InvestorReportStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	investorReportStructErr := stub.PutState(irID, []byte(investorReportStructBytes))
	if investorReportStructErr != nil {
		fmt.Println("Couldn't save investor report Characterestic data to ledger", investorReportStructErr)
		return "", fmt.Errorf("Couldn't save investor report Characterestic data to ledger")
	}

	indexName := "investorReportMonthIndex"
	irMonthYearKey, err := stub.CreateCompositeKey(indexName, []string{irDealID, irMonth, irYear, irID})
	if err != nil {
		fmt.Println("Could not index composite key for month, year and deal id", err)
		return "", fmt.Errorf("Could not index composite key for month and year map")
	}

	value := []byte{0x00}
	stub.PutState(irMonthYearKey, value)

	var irEvent = "{eventType: 'InvestorReportCC', description:" + InvestorReportStruct.IRMonth + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(irEvent))
	if err != nil {
		return "", fmt.Errorf("Couldn't set event hub")
	}
	fmt.Println("Successfully saved InvestorReport")
	return InvestorReportStruct.IRHash, nil
}

func main() {
	err := shim.Start(new(InvestorReport))
	{
		fmt.Printf("Error starting InvestorReport chaincode: %s", err)
	}
}
