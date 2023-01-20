package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type CreatePoolAttr struct {
	PoolAttrID    string `json:"PoolAttrID"`
	PoolAttrName  string `json:"PoolAttrName"`
	PoolAttrValue string `json:"PoolAttrValue"`
	Status        string `json:"Status"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *CreatePoolAttr) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *CreatePoolAttr) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "CreatePoolAttrs" {

		result, err = CreatePoolAttrs(stub, args)
	}
	if fn == "GetpoolAttr" {

		result, err = GetpoolAttr(stub, args)
	}
	if fn == "GetAllpoolAttrs" {

		result, err = GetAllpoolAttrs(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetpoolAttr(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetpoolAttr")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing PoolAttributeID")
	}

	var PoolID = args[0]
	value, err := stub.GetState(PoolID)
	if err != nil {
		fmt.Println("Could not get Pool Attribute ID  with id "+PoolID+" from ledger", err)
		return "", errors.New("Missing PoolID")
	}
	return string(value), nil
}

func GetAllpoolAttrs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetAllpools")

	GetAllPoolAttrsItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Pools")
	}
	defer GetAllPoolAttrsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllPoolAttrsItr.HasNext() {
		queryResponseValue, err := GetAllPoolAttrsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Pool Data")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponseValue.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponseValue.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- GetAllPoolAttrs queryResult:\n%s\n", buffer.String())

	return buffer.String(), nil
}

func CreatePoolAttrs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CreatePoolAttr")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var poolAttrID = args[0]
	fmt.Println("the Pool is" + poolAttrID)
	var poolAttrName = args[1]
	fmt.Println("the poolAttrName is" + poolAttrName)
	var PoolAttrValue = args[2]
	fmt.Println("the PoolAttrValue is" + PoolAttrValue)

	//assigning to struct the variables
	createPoolAttrStruct := CreatePoolAttr{
		PoolAttrID:    poolAttrID,
		PoolAttrName:  poolAttrName,
		PoolAttrValue: PoolAttrValue,
		Status:        "active",
	}

	fmt.Println("the struct values are PoolAttrID  " + createPoolAttrStruct.PoolAttrID + "PoolAttrName" + createPoolAttrStruct.PoolAttrName + "PoolAttrValue" + createPoolAttrStruct.PoolAttrValue)

	poolStructBytes, err := json.Marshal(createPoolAttrStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	poolerr := stub.PutState(poolAttrID, []byte(poolStructBytes))
	if poolerr != nil {
		fmt.Println("Could not save poolCharacterestic data to ledger", poolerr)
		return "", fmt.Errorf("Could not save poolCharacterestic data to ledger")
	}

	var poolEvent = "{eventType: 'CreatePoolAttr', description:" + createPoolAttrStruct.PoolAttrID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(poolEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved CreatePoolAttr")
	return args[0], nil

}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(CreatePoolAttr))
	{
		fmt.Printf("Error starting CreatePoolAttr chaincode: %s", err)
	}
}
