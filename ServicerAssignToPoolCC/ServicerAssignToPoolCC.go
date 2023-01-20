package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ServicerAssignToPool struct {
	PoolID         string `json:"PoolID"`
	ServicerID     string `json:"ServicerID"`
	ApprovedStatus string `json:"ApprovedStatus"`
	ApprovedDate   string `json:"ApprovedDate"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *ServicerAssignToPool) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *ServicerAssignToPool) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SaveApproval" {

		result, err = SaveApproval(stub, args)
	}
	if fn == "GetPoolService" {

		result, err = GetPoolService(stub, args)
	}
	if fn == "GetPoolsAssignedToServicer" {

		result, err = GetPoolsAssignedToServicer(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetPoolService(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering getPoolService")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing PoolID")
	}

	var Pool_ID = args[0]
	value, err := stub.GetState(Pool_ID)
	if err != nil {
		fmt.Println("Could not get Pool_ID  with id "+Pool_ID+" from ledger", err)
		return "", errors.New("Missing Pool_ID")
	}
	return string(value), nil
}

func SaveApproval(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering saveApproval")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var pool_ID = args[0]
	fmt.Println("the Pool is" + pool_ID)
	var servicer_id = args[1]
	fmt.Println("the servicer_id is" + servicer_id)
	var appov_sts = args[2]
	fmt.Println("the appov_sts is" + appov_sts)
	var approv_date = args[3]
	fmt.Println("the approv_date is" + approv_date)

	//assigning to struct the variables
	ServicerAssignToPoolStruct := ServicerAssignToPool{
		PoolID:         pool_ID,
		ServicerID:     servicer_id,
		ApprovedStatus: appov_sts,
		ApprovedDate:   approv_date,
	}

	fmt.Println("the struct values are PoolID  " + ServicerAssignToPoolStruct.PoolID + "servicerID" + ServicerAssignToPoolStruct.ServicerID + "approvedStatus" + ServicerAssignToPoolStruct.ApprovedStatus)

	poolStructBytes, err := json.Marshal(ServicerAssignToPoolStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	poolerr := stub.PutState(pool_ID, []byte(poolStructBytes))
	if poolerr != nil {
		fmt.Println("Could not save ServicerAssignToPool data to ledger", poolerr)
		return "", fmt.Errorf("Could not save ServicerAssignToPool data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName := "servicerPoolIndex"
	poolServicerKey, err := stub.CreateCompositeKey(indexName, []string{ServicerAssignToPoolStruct.ServicerID, ServicerAssignToPoolStruct.PoolID})
	if err != nil {
		fmt.Println("Could not index composite key for pool servicer map", err)
		return "", fmt.Errorf("Could not index composite key for pool servicer")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(poolServicerKey, value)

	var poolEvent = "{eventType: 'saveApproval', description:" + ServicerAssignToPoolStruct.PoolID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(poolEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved ServicerAssignToPool")
	return args[0], nil

}

//get pools assigned to Servicer
func GetPoolsAssignedToServicer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPoolsAssignedToServicer")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing servicerId")
	}
	//form client pool id is the key which is created by a custom logic
	var servicerID = args[0]
	servicerPoolsItr, err := stub.GetStateByPartialCompositeKey("servicerPoolIndex", []string{servicerID})
	if err != nil {
		return "", fmt.Errorf("Could not get pools for Servicer")
	}
	defer servicerPoolsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"PoolsForServicer\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for servicerPoolsItr.HasNext() {
		queryResponse, err := servicerPoolsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next pool for Servicer")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the servicerId and poolId from poolAssignApprovalIndex composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		returnedServicerID := compositeKeyParts[0]
		returnedPoolID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s servicerId:%s poolId:%s\n", objectType, returnedServicerID, returnedPoolID)

		value, err := stub.GetState(returnedPoolID)
		if err != nil {
			fmt.Println("Could not get servicer for id "+returnedPoolID+" from ledger", err)
			return "", errors.New("Missing PoolID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	fmt.Printf("- GetPoolsAssignedToServicer queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(ServicerAssignToPool))
	{
		fmt.Printf("Error starting ServicerAssignToPool chaincode: %s", err)
	}
}
