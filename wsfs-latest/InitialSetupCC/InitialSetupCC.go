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

type InitialSetup struct {
	ISID                      string `json:"isId"`
	ISDealId                  string `json:"isDealId"`
	ISMaturityDate            string `json:"isMaturityDate"`
	ISNodInitialAccrualPeriod string `json:"isNodInitialAccrualPeriod"`
	ISNodAccrualPeriod        string `json:"isNodAccrualPeriod"`
	ISUpdatedBy               string `json:"isUpdatedBy"`
	ISUpdationDate            string `json:"isUpdationDate"`
	ISReserveFundBalance      string `json:"isReserveFundBalance"`
	ISCertificateBalance      string `json:"isCertificateBalance"`
	ISCutoffLoanCount 				string `json:"isCutoffLoanCount"`
	ISCutoffLoanBalance 			string `json:"isCutoffLoanBalance"`
	ISCutoffAverageBalance 		string `json:"isCutoffAverageBalance"`
	ISCutoffWtdAvgRate 				string `json:"isCutoffWtdAvgRate"`
	ISCutoffWtdAvgTerm 				string `json:"isCutoffWtdAvgTerm"`
}

func (t *InitialSetup) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *InitialSetup) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetInitialSetupByID" {

		result, err = GetInitialSetupByID(stub, args)
	}
	if fn == "GetInitialSetupByDealID" {

		result, err = GetInitialSetupByDealID(stub, args)
	}
	if fn == "SaveInitialSetup" {

		result, err = SaveInitialSetup(stub, args)
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

func GetInitialSetupByID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInitialSetupByID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing ISID")
	}

	var isId = args[0]
	value, err := stub.GetState(isId)
	if err != nil {
		fmt.Println("Couldn't get Initial Setup with id "+isId+" from ledger", err)
		return "", errors.New("Missing isId")
	}

	return string(value), nil
}

func GetInitialSetupByDealID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetInitialSetupByDealID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing dealId")
	}

	var dealId = args[0]

	getISByDealIdItr, err := stub.GetStateByPartialCompositeKey("initialSetupDealIndex", []string{dealId})
	if err != nil {
		return "", fmt.Errorf("Could not get Initial Setup for this deal")
	}
	defer getISByDealIdItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	// buffer.WriteString("{\"DTOInitialSetup\":")
	bArrayMemberAlreadyWritten := false
	for getISByDealIdItr.HasNext() {
		queryResponse, err := getISByDealIdItr.Next()
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
		returnedISDealId := compositeKeyParts[0]
		returnedISID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s isID:%s isDealId:%s\n", objectType, returnedISID, returnedISDealId)

		value, err := stub.GetState(returnedISID)
		if err != nil {
			fmt.Println("Couldn't get initial setup for id "+returnedISID+" from ledger", err)
			return "", errors.New("Missing isID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
		break
	}
	// buffer.WriteString("}")

	return buffer.String(), nil
}

func SaveInitialSetup(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveInitialSetup")

	if len(args) < 6 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var isID = args[0]
	fmt.Println("the ISID is" + isID)
	var isDealId = args[1]
	fmt.Println("the Initial Setup is" + isDealId)
	var isMaturityDate = args[2]
	fmt.Println("the Maturity Date is" + isMaturityDate)
	var isNodInitialAccrualPeriod = args[3]
	fmt.Println("the Initial Accrual Period is" + isNodInitialAccrualPeriod)
	var isNodAccrualPeriod = args[4]
	fmt.Println("the Accrual Period is" + isNodAccrualPeriod)
	var isUpdatedBy = args[5]
	fmt.Println("the User ID is" + isUpdatedBy)
	var isUpdationDate = args[6]
	fmt.Println("the Updation Date is" + isUpdationDate)
	var isReserveFundBalance = args[7]
	fmt.Println("the reserv fund balance is " + isReserveFundBalance)
	var isCertificateBalance = args[8]
	fmt.Println("the certificate balance is " + isCertificateBalance)
	var isCutoffLoanCount = args[9]
	fmt.Println("the Cutoff Loan Count is" + isCutoffLoanCount)
	var isCutoffLoanBalance = args[10]
	fmt.Println("the Cutoff Loan Balance is" + isCutoffLoanBalance)
	var isCutoffAverageBalance = args[11]
	fmt.Println("the Cutoff Avg Balance is" + isCutoffAverageBalance)
	var isCutoffWtdAvgRate = args[12]
	fmt.Println("the Cutoff Wtd Avg Rate is " + isCutoffWtdAvgRate)
	var isCutoffWtdAvgTerm = args[13]
	fmt.Println("the Cutoff Wtd Avg Term is " + isCutoffWtdAvgTerm)

	InitialSetupStruct := InitialSetup{
		ISID:                      isID,
		ISDealId:                  isDealId,
		ISMaturityDate:            isMaturityDate,
		ISNodInitialAccrualPeriod: isNodInitialAccrualPeriod,
		ISNodAccrualPeriod:        isNodAccrualPeriod,
		ISUpdatedBy:               isUpdatedBy,
		ISUpdationDate:            isUpdationDate,
		ISReserveFundBalance:      isReserveFundBalance,
		ISCertificateBalance:      isCertificateBalance,
		ISCutoffLoanCount: 				 isCutoffLoanCount,
		ISCutoffLoanBalance: 			 isCutoffLoanBalance,
		ISCutoffAverageBalance:		 isCutoffAverageBalance,
		ISCutoffWtdAvgRate:				 isCutoffWtdAvgRate,
		ISCutoffWtdAvgTerm: 			 isCutoffWtdAvgTerm,
	}

	fmt.Println("the struct values are ISID:" + InitialSetupStruct.ISID + " ISDealId:" + InitialSetupStruct.ISDealId +
		" ISMaturityDate:" + InitialSetupStruct.ISMaturityDate)

	initialSetupStructBytes, err := json.Marshal(InitialSetupStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	initialSetupStructErr := stub.PutState(isID, []byte(initialSetupStructBytes))
	if initialSetupStructErr != nil {
		fmt.Println("Couldn't save initial setup Characterestic data to ledger", initialSetupStructErr)
		return "", fmt.Errorf("Couldn't save initial setup Characterestic data to ledger")
	}

	indexName := "initialSetupDealIndex"
	isDealIdKey, err := stub.CreateCompositeKey(indexName, []string{isDealId, isID})
	if err != nil {
		fmt.Println("Could not index composite key for dealId", err)
		return "", fmt.Errorf("Could not index composite key for dealId")
	}

	value := []byte{0x00}
	stub.PutState(isDealIdKey, value)

	var isEvent = "{eventType: 'InitialSetupCC', description:" + InitialSetupStruct.ISDealId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(isEvent))
	if err != nil {
		return "", fmt.Errorf("Couldn't set event hub")
	}
	fmt.Println("Successfully saved InitialSetup")
	return "", nil
}

func main() {
	err := shim.Start(new(InitialSetup))
	{
		fmt.Printf("Error starting InitialSetup chaincode: %s", err)
	}
}
