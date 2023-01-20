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

type Tranche struct {
  TRID string `json:"trID"`
  TRDealId string `json:"trDealId"`
	TRNote string `json:"trNote"`
	TRCusip string `json:"trCusip"`
	TROriginalBalance string `json:"trOriginalBalance"`
	TRInterestRate string `json:"trInterestRate"`
	TRUpdatedBy string `json:"trUpdatedBy"`
	TRUpdationDate string `json:"trUpdationDate"`
}

func (t *Tranche) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *Tranche) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetTrancheByID" {

		result, err = GetTrancheByID(stub, args)
	}
	if fn == "GetTrancheByDealID" {

		result, err = GetTrancheByDealID(stub, args)
	}
	if fn == "SaveTranche" {

		result, err = SaveTranche(stub, args)
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

func GetTrancheByID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheByID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing ISID")
	}

	var trancheId = args[0]
	value, err := stub.GetState(trancheId)
	if err != nil {
		fmt.Println("Couldn't get Tranche with id "+trancheId+" from ledger", err)
		return "", errors.New("Missing trancheId")
	}

	return string(value), nil
}

func GetTrancheByDealID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetTrancheByDealID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing dealId")
	}

	var dealId = args[0]

	getTrancheByDealIdItr, err := stub.GetStateByPartialCompositeKey("trancheDealIndex", []string{dealId})
	if err != nil {
		return "", fmt.Errorf("Could not get Tranche for this deal")
	}
	defer getTrancheByDealIdItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"DTOTranche\":")
	bArrayMemberAlreadyWritten := false
	buffer.WriteString("[")
	for getTrancheByDealIdItr.HasNext() {
		queryResponse, err := getTrancheByDealIdItr.Next()
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
		returnedTRDealId := compositeKeyParts[0]
		returnedTRID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s trID:%s trDealId:%s\n", objectType, returnedTRID, returnedTRDealId)

		value, err := stub.GetState(returnedTRID)
		if err != nil {
			fmt.Println("Couldn't get tranche for id "+returnedTRID+" from ledger", err)
			return "", errors.New("Missing trID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func SaveTranche(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveTranche")

	if len(args) < 6 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var trID = args[0]
	fmt.Println("the trID is" + trID)
	var trDealId = args[1]
	fmt.Println("the Deal ID is" + trDealId)
	var trNote = args[2]
	fmt.Println("the Note is" + trNote)
  	var trCusip = args[3]
	fmt.Println("the Cusip is" + trCusip)
  	var trOriginalBalance = args[4]
	fmt.Println("the Original Balance is" + trOriginalBalance)
  	var trInterestRate = args[5]
	fmt.Println("the Interest Rate is" + trInterestRate)
  	var trUpdatedBy = args[6]
	fmt.Println("the User ID is" + trUpdatedBy)
  	var trUpdationDate = args[7]
	fmt.Println("the Updation Date is" + trUpdationDate)

	TrancheStruct := Tranche{
    	TRID: trID,
    	TRDealId: trDealId,
  	TRNote: trNote,
  	TRCusip: trCusip,
  	TROriginalBalance: trOriginalBalance,
  	TRInterestRate: trInterestRate,
  	TRUpdatedBy: trUpdatedBy,
  	TRUpdationDate: trUpdationDate,
	}

	fmt.Println("the struct values are TRID:" + TrancheStruct.TRID + " TRDealId:" + TrancheStruct.TRDealId +
		" TRNote:" + TrancheStruct.TRNote + " TRCusip:"+ TrancheStruct.TRCusip + " TROriginalBalance:" + TrancheStruct.TROriginalBalance)

  trancheStructBytes, err := json.Marshal(TrancheStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	trancheStructErr := stub.PutState(trID, []byte(trancheStructBytes))
	if trancheStructErr != nil {
		fmt.Println("Couldn't save tranche Characterestic data to ledger", trancheStructErr)
		return "", fmt.Errorf("Couldn't save tranche Characterestic data to ledger")
	}

	indexName := "trancheDealIndex"
	trDealIdKey, err := stub.CreateCompositeKey(indexName, []string{trDealId, trID})
	if err != nil {
		fmt.Println("Could not index composite key for dealId", err)
		return "", fmt.Errorf("Could not index composite key for dealId")
	}

	value := []byte{0x00}
	stub.PutState(trDealIdKey, value)

	var trEvent = "{eventType: 'TrancheCC', description:" + TrancheStruct.TRDealId+ "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(trEvent))
	if err != nil {
		return "", fmt.Errorf("Couldn't set event hub")
	}
	fmt.Println("Successfully saved InitialSetup")
	return "", nil
}

func main() {
	err := shim.Start(new(Tranche))
	{
		fmt.Printf("Error starting InitialSetup chaincode: %s", err)
	}
}
