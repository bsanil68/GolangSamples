package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/*type DLoanID struct {
	ID string `json:"ID"`
}*/

type DLoanMapCC struct {
	PoolID string `json:"PoolID"`
	//LoanID []DLoanID `json:"LoanID"`
	LoanID string `json:"LoanID"`
}

func (t *DLoanMapCC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *DLoanMapCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetLoanMapByPoolID" {

		result, err = GetLoanMapByPoolID(stub, args)
	}

	if fn == "GetAllLoanMapDetails" {
		result, err = GetAllLoanMapDetails(stub)
	}
	if fn == "SaveLoanMapDetails" {

		err = SaveLoanMapDetails(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetLoanMapByPoolID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoanMapByPoolID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing Pool Id")
	}

	var PoolID = args[0]
	value, err := stub.GetState(PoolID)
	if err != nil {
		fmt.Println("Couldn't get Pool Id  with id "+PoolID+" from ledger", err)
		return "", errors.New("Missing Pool Id")
	}

	return string(value), nil
}

func GetAllLoanMapDetails(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllLoanMapDetails")

	getAllLoanMapItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all LoanMap details")
	}
	defer getAllLoanMapItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"GetAllLoanMapDetails\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for getAllLoanMapItr.HasNext() {
		queryResponseValue, err := getAllLoanMapItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Data")
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

func SaveLoanMapDetails(stub shim.ChaincodeStubInterface, args []string) error {
	fmt.Println("Entering SaveLoanMapDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	poolID := args[0]
	//loanIdsStr := args[1] //"l1,l2,l3"
	loanID := args[1] //"l1,l2,l3"
	//fmt.Println("loanIdsStr")
	//loanID := strings.Split(loanIdsStr, ",")

	//item := make([]DLoanID, len(loanID))
	/*var k = 0
	for i, j := range args {
		if i < 1 {
			poolID = j
			fmt.Println("the ID is" + poolID)
		} else {
			loanID[k] = j
			fmt.Println("the Iterator before increment is" + string(k))
			fmt.Println("the LoanID is" + loanID[k])
			k++
			fmt.Println("the Iterator after increment is" + string(k))
		}

	}*/
	//fmt.Println(len(loanID))
	//int k=0
	/*for loopcount := 0; loopcount < len(loanID); loopcount++ {
		item[loopcount].ID = loanID[loopcount]
	}*/
	//for itr := 0; itr < len(loanID); itr++ {
	LoanMapCCStruct := DLoanMapCC{
		PoolID: poolID,
		LoanID: loanID,
	}

	//fmt.Println(LoanMapCCStruct.PoolID + " " + LoanMapCCStruct.LoanID + " " + itr)
	LoanMapStructBytes, err := json.Marshal(LoanMapCCStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return fmt.Errorf("Couldn't marshal data from struct")

	}
	LoanMapErr := stub.PutState(poolID, []byte(LoanMapStructBytes))
	if LoanMapErr != nil {
		fmt.Println("Couldn't save LoanMapCharacterestic data to ledger", LoanMapErr)
		return fmt.Errorf("Couldn't save LoanMapCharacterestic data to ledger")
	}

	fmt.Println("Successfully saved ")

	//}
	return nil
}

func main() {
	err := shim.Start(new(DLoanMapCC))
	{
		fmt.Printf("Error starting LoanMapCC chaincode: %s", err)
	}
}
