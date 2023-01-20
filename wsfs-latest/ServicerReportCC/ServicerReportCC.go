package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type ServicerReport struct {
	SRID           string `json:"srId"`
	SRKey          string `json:"srKey"`
	SRValue        string `json:"srValue"`
	SRMonth        string `json:"srMonth"`
	SRYear         string `json:"srYear"`
	SRUpdatedBy    string `json:"srUpdatedBy"`
	SRUpdationDate string `json:"srUpdationDate"`
	SRSeqNum       string `json:"srSeqNum"`
	SRDealID       string `json:"srDealId"`
}

func (t *ServicerReport) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *ServicerReport) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if function == "SaveServicerReportAttribute" {
		// Deletes an entity from its state
		return t.SaveServicerReportAttribute(stub, args)
	}

	if function == "GetServicerReportAttributeBySRID" {
		// queries an entity state
		return t.GetServicerReportAttributeBySRID(stub, args)
	}

	if function == "GetServicerReportAttributeBySRKey" {
		// queries an entity state
		return t.GetServicerReportAttributeBySRKey(stub, args)
	}

	if function == "GetServicerReportAttributeByMonthAndYear" {
		// queries an entity state
		return t.GetServicerReportAttributeByMonthAndYear(stub, args)
	}

	if function == "GetAllServicerReportAttributes" {
		// queries an entity state
		return t.GetAllServicerReportAttributes(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func (t *ServicerReport) SaveServicerReportAttribute(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) < 7 {
		return shim.Error("Incorrect arguments. Expecting 8 arguments")
	}
	// json args input musts match  the case and  spelling exactly
	// get all the  arguments
	var sRID = args[0]
	fmt.Println("the sRID is" + sRID)
	var sRKey = args[1]
	fmt.Println("the sRKey is" + sRKey)
	var sRValue = args[2]
	fmt.Println("the sRValue is" + sRValue)
	var sRMonth = args[3]
	fmt.Println("the sRMonth is" + sRMonth)
	var sRYear = args[4]
	fmt.Println("the sRYear is" + sRYear)
	var sRUpdatedBy = args[5]
	fmt.Println("the sRUpdatedBy is" + sRUpdatedBy)
	var sRUpdationDate = args[6]
	fmt.Println("the sRUpdationDate is" + sRUpdationDate)
	var sRSeqNum = args[7]
	fmt.Println("the sRSeqNum is" + sRSeqNum)
	var sRDealID = args[8]
	fmt.Println("the deal id is " + sRDealID)

	//assigning to struct the variables
	servicerReportStruct := ServicerReport{
		SRID:           sRID,
		SRKey:          sRKey,
		SRValue:        sRValue,
		SRMonth:        sRMonth,
		SRYear:         sRYear,
		SRUpdatedBy:    sRUpdatedBy,
		SRUpdationDate: sRUpdationDate,
		SRSeqNum:       sRSeqNum,
		SRDealID:       sRDealID,
	}

	//  show that  the struct  has  values false

	fmt.Println("the struct values are SRID  " + servicerReportStruct.SRID + "SRValue" + servicerReportStruct.SRValue + "SRUpdatedBy" + servicerReportStruct.SRUpdatedBy)

	servicerReportStructBytes, err := json.Marshal(servicerReportStruct)
	if err != nil {
		return shim.Error("Unable to unmarshal data from struct")

	}
	logerr := stub.PutState(sRID, []byte(servicerReportStructBytes))
	if logerr != nil {
		return shim.Error(logerr.Error())
	}

	var servicerEvent = "{eventType: 'SaveServicerReportAttribute', description:" + servicerReportStruct.SRID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(servicerEvent))
	if err != nil {
		return shim.Error("Could not set event hub")
	}

	sRKeyIndex := "sRKeyIndex"
	SRKeyIDCompKey, err := stub.CreateCompositeKey(sRKeyIndex, []string{servicerReportStruct.SRKey, servicerReportStruct.SRID})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(SRKeyIDCompKey, value)

	MonthAndYearIndex := "MonthAndYearIndex"
	monthAndYearCompKey, err := stub.CreateCompositeKey(MonthAndYearIndex, []string{servicerReportStruct.SRDealID, servicerReportStruct.SRMonth, servicerReportStruct.SRYear, servicerReportStruct.SRID})
	if err != nil {
		return shim.Error(err.Error())
	}

	value1 := []byte{0x00}
	stub.PutState(monthAndYearCompKey, value1)

	fmt.Println("Successfully saved ServicerReport")

	return shim.Success(nil)
}

func (t *ServicerReport) GetAllServicerReportAttributes(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//Get all Users from the ledger
	GetAllServicerReportItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error("Could not get all Schdulars")
	}
	defer GetAllServicerReportItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"DTOServicerReport\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllServicerReportItr.HasNext() {
		queryResponseValue, err := GetAllServicerReportItr.Next()
		if err != nil {
			return shim.Error("Could not get next Users")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponseValue.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())

}

func (t *ServicerReport) GetServicerReportAttributeBySRKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var sRKey string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting sRKey to query")
	}

	sRKey = args[0]

	//Get all stages by schID
	sRKeyItr, err := stub.GetStateByPartialCompositeKey("sRKeyIndex", []string{sRKey})

	if err != nil {
		return shim.Error("Could not get composite key for this sRKey")
	}
	defer sRKeyItr.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	//buffer.WriteString("{\"DTOServicerReport\":")
	//buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for sRKeyItr.HasNext() {
		queryResponse, err := sRKeyItr.Next()
		if err != nil {
			return shim.Error("Could not get next stage")
		}

		// Add a comma before arrat members, supress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the customerId from the customerId~loanId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error("Could not get Composite key parts")
		}

		returnedSRKey := compositeKeyParts[0]
		returnedSRID := compositeKeyParts[1]

		fmt.Printf("- found a from index:%s SRKey:%s SRID:%s", objectType, returnedSRKey, returnedSRID)

		value, err := stub.GetState(returnedSRID)
		if err != nil {
			return shim.Error("Missing SRID in ledger")
		}

		// Returned Record is a JSON object. So write as is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	//buffer.WriteString("]")
	//	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

func (t *ServicerReport) GetServicerReportAttributeByMonthAndYear(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var month string
	var year string
	var dealID string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting month to query")
	}

	dealID = args[0]
	month = args[1]
	year = args[2]

	//Get all stages by schID
	monthAndYearItr, err := stub.GetStateByPartialCompositeKey("MonthAndYearIndex", []string{dealID, month, year})

	if err != nil {
		return shim.Error("Could not get composite key for this month")
	}
	defer monthAndYearItr.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("{\"DTOServicerReport\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for monthAndYearItr.HasNext() {
		queryResponse, err := monthAndYearItr.Next()
		if err != nil {
			return shim.Error("Could not get next stage")
		}

		// Add a comma before arrat members, supress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the customerId from the customerId~loanId composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error("Could not get Composite key parts")
		}
		returnedDealID := compositeKeyParts[0]
		returnedMonth := compositeKeyParts[1]
		returnedYear := compositeKeyParts[2]
		returnedSRID := compositeKeyParts[3]

		fmt.Printf("- found a from index:%s month:%s year:%s SRID:%s DealID:%s", objectType, returnedMonth, returnedYear, returnedSRID, returnedDealID)

		value, err := stub.GetState(returnedSRID)
		if err != nil {
			return shim.Error("Missing SRID in ledger")
		}

		// Returned Record is a JSON object. So write as is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

// Query callback representing the query of a chaincode
func (t *ServicerReport) GetServicerReportAttributeBySRID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Entering GetServicerReportAttributeBySRID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return shim.Error("Missing Servicer Report Id")
	}

	var buffer bytes.Buffer

	//buffer.WriteString("{\"ServicerReport\":")

	//  by a  custom logic
	var SRID = args[0]
	servicerReportStructBytes, err := stub.GetState(SRID)

	if err != nil {
		fmt.Println("Could not get Servicer report with id "+SRID+" from ledger", err)
		return shim.Error("Missing Servicer Report Id")
	}
	buffer.WriteString(string(servicerReportStructBytes))
	// buffer.WriteString("}")

	return shim.Success(buffer.Bytes())

}

func main() {
	err := shim.Start(new(ServicerReport))
	if err != nil {

	}
}
