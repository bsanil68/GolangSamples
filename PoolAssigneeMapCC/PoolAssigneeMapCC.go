package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type PoolAssignee struct {
	PoolID         string `json:"PoolID"`
	AssigneeID     string `json:"AssigneeID"`
	ApprovedDate   string `json:"ApprovedDate"`
	ApproverType   string `json:"ApproverType"`
	ApproverID     string `json:"ApproverID"`
	ApprovedStatus string `json:"ApprovedStatus"`
}

// we do not init  any data   when  chaincode  is called
func (t *PoolAssignee) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)

}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *PoolAssignee) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SavePoolAssigneeMap" {

		result, err = SavePoolAssigneeMap(stub, args)
	}
	if fn == "GetPoolAssigneeMap" {

		result, err = GetPoolAssigneeMap(stub, args)
	}
	if fn == "GetAllPoolAssigneeMaps" {

		result, err = GetAllPoolAssigneeMaps(stub)
	}
	if fn == "GetAllPoolsForAssignee" {

		result, err = GetAllPoolsForAssignee(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

//assign pool to Assignee
func SavePoolAssigneeMap(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SavePoolAssigneeMap")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var poolID = args[0]
	fmt.Println("the Pool is" + poolID)
	var assigneeID = args[1]
	fmt.Println("the assigner_id is" + assigneeID)
	var approvDate = args[2]
	fmt.Println("the approv_date is" + approvDate)
	var approvType = args[3]
	fmt.Println("the approv_type is" + approvType)
	var approvId = args[4]
	fmt.Println("the approv_id is" + approvId)

	var status = args[5]
	fmt.Println("the status is" + approvId)

	//assigning to struct the variables
	PoolAssigneeStruct := PoolAssignee{
		PoolID:         poolID,
		AssigneeID:     assigneeID,
		ApprovedDate:   approvDate,
		ApproverType:   approvType,
		ApproverID:     approvId,
		ApprovedStatus: status,
	}

	fmt.Println("the struct values are PoolID  " + PoolAssigneeStruct.PoolID + "assigneeID" + PoolAssigneeStruct.AssigneeID + "approvedDate" + PoolAssigneeStruct.ApprovedDate)

	poolStructBytes, err := json.Marshal(PoolAssigneeStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}

	poolerr := stub.PutState(poolID, []byte(poolStructBytes))
	if poolerr != nil {
		fmt.Println("Could not save PoolAssignApproval data to ledger", poolerr)
		return "", fmt.Errorf("Could not save PoolAssignApproval data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "poolAssigneeIndex"
	assigneeKey, err := stub.CreateCompositeKey(indexName, []string{PoolAssigneeStruct.AssigneeID, PoolAssigneeStruct.PoolID})
	if err != nil {
		fmt.Println("Could not index composite key for pool and assignee", err)
		return "", fmt.Errorf("Could not index composite key for pool assignee map")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(assigneeKey, value)

	var poolEvent = "{eventType: 'saveAssignementApproval', description:" + PoolAssigneeStruct.PoolID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(poolEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved PoolAssignApproval")
	return args[0], nil

}

func GetPoolAssigneeMap(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPoolAssigneeMap")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing assigneeId")
	}

	var buffer bytes.Buffer

	//	buffer.WriteString("{\"PoolAssigneeMap\":")

	var PoolID = args[0]
	value, err := stub.GetState(PoolID)
	if err != nil {
		fmt.Println("Could not get assigner for id "+PoolID+" from ledger", err)
		return "", errors.New("Missing PoolID")
	}
	buffer.WriteString(string(value))
	//	buffer.WriteString("}")
	return buffer.String(), nil
}

func GetAllPoolAssigneeMaps(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllPoolAssigneeMaps")

	GetAllAssineeMapsItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Pool Assignee Maps")
	}
	defer GetAllAssineeMapsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"PoolAssigneMapList\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllAssineeMapsItr.HasNext() {
		queryResponseValue, err := GetAllAssineeMapsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Pool Assinee Data")
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
	//fmt.Printf("- GetAllPoolAssigneeMaps queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

func GetAllPoolsForAssignee(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetAllPoolsForAssignee")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing document")
	}

	assigneeID := args[0]

	GetAllPoolsForAssigneeItr, err := stub.GetStateByPartialCompositeKey("poolAssigneeIndex", []string{assigneeID})
	if err != nil {
		return "", fmt.Errorf("Could not get all Pool Assignee Maps")
	}
	defer GetAllPoolsForAssigneeItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"PoolsForAssignee\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllPoolsForAssigneeItr.HasNext() {
		queryResponse, err := GetAllPoolsForAssigneeItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Pool Assinee Data")
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
		returnedAssigneeID := compositeKeyParts[0]
		returnedPoolID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s loanid:%s poolid:%s\n", objectType, returnedAssigneeID, returnedPoolID)

		value, err := stub.GetState(returnedPoolID)
		if err != nil {
			fmt.Println("Could not get assigner for id "+returnedPoolID+" from ledger", err)
			return "", errors.New("Missing PoolID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")
	//fmt.Printf("- GetAllPoolAssigneeMaps queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(PoolAssignee))
	{
		fmt.Printf("Error starting poolassigneeMapCC chaincode: %s", err)
	}
}
