package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type CreatePool struct {
	PoolID          string `json:"PoolID"`
	PoolName        string `json:"PoolName"`
	PoolDesc        string `json:"PoolDesc"`
	Poolowner       string `json:"Poolowner"`
	PoolStartedDate string `json:"PoolStartedDate"`
	PoolExpiryDate  string `json:"PoolExpiryDate "`
	PoolStatus      string `json:"PoolStatus"`
	ApprovalDate    string `json:"ApprovalDate"`
	ApproverID      string `json:"ApproverID"`
	NoOfAssets      string `json:"NoOfAssets"`
	Status          string `json:"Status"`
	PoolHash        string `json:"PoolHash"`
	//LoanHash        string `json:"LoanHash"`
	PoolCreatedDate     string `json:"PoolCreatedDate"`
	WarehouseLenderName string `json:"WarehouseLenderName"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *CreatePool) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *CreatePool) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "Createpool" {

		result, err = Createpool(stub, args)
	}
	if fn == "Getpool" {

		result, err = Getpool(stub, args)
	}
	if fn == "GetPoolsByDate" {

		result, err = GetPoolsByDate(stub, args)
	}
	if fn == "GetPoolsByWarehouseLenderName" {

		result, err = GetPoolsByWarehouseLenderName(stub, args)
	}
	if fn == "GetAllpools" {

		result, err = GetAllpools(stub)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func Getpool(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering Getpool")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing PoolID")
	}

	var buffer bytes.Buffer

	//buffer.WriteString("{")

	var PoolID = args[0]
	value, err := stub.GetState(PoolID)

	if err != nil {
		fmt.Println("Could not get PoolID  with id "+PoolID+" from ledger", err)
		return "", errors.New("Missing PoolID")
	}
	buffer.WriteString(string(value))
	//buffer.WriteString("}")
	return buffer.String(), nil
}

func GetPoolsByDate(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPoolsByDate")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing createdDate")
	}

	var createdDate = args[0]

	GetPoolsByDateItr, err := stub.GetStateByPartialCompositeKey("createDateIndex", []string{createdDate})
	if err != nil {
		return "", fmt.Errorf("Could not get pools for this date")
	}
	defer GetPoolsByDateItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetPoolsByDateItr.HasNext() {
		queryResponse, err := GetPoolsByDateItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next pool Data")
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
		returnedcreateDate := compositeKeyParts[0]
		returnedpoolID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s createDate:%s poolid:%s\n", objectType, returnedcreateDate, returnedpoolID)

		value, err := stub.GetState(returnedpoolID)
		if err != nil {
			fmt.Println("Couldn't get pools for id "+returnedpoolID+" from ledger", err)
			return "", errors.New("Missing poolID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetPoolsByWarehouseLenderName(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetPoolsByWarehouseLenderName")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing warehouse lender name")
	}

	var warehouseLenderName = args[0]

	GetPoolsByLenderNameItr, err := stub.GetStateByPartialCompositeKey("LenderNameIndex", []string{warehouseLenderName})
	if err != nil {
		return "", fmt.Errorf("Could not get pools for this lender name")
	}
	defer GetPoolsByLenderNameItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetPoolsByLenderNameItr.HasNext() {
		queryResponse, err := GetPoolsByLenderNameItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next pool Data")
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
		returnedlenderName := compositeKeyParts[0]
		returnedpoolID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s LenderName:%s poolid:%s\n", objectType, returnedlenderName, returnedpoolID)

		value, err := stub.GetState(returnedpoolID)
		if err != nil {
			fmt.Println("Couldn't get pools for id "+returnedpoolID+" from ledger", err)
			return "", errors.New("Missing poolID")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetAllpools(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllpools")

	GetAllPoolsItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Pools")
	}
	defer GetAllPoolsItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllPoolsItr.HasNext() {
		queryResponseValue, err := GetAllPoolsItr.Next()
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

func Createpool(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering createpool")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var poolID = args[0]
	fmt.Println("the Pool is" + poolID)
	var poolName = args[1]
	fmt.Println("the poolName is" + poolName)
	var poolDesc = args[2]
	fmt.Println("the poolDesc is" + poolDesc)
	var poolOwner = args[3]
	fmt.Println("the poolOwner is" + poolOwner)
	var poolStartedDate = args[4]
	fmt.Println("the PoolStartedDate is" + poolStartedDate)
	var poolExpiryDate = args[5]
	fmt.Println("the PoolExpiryDate is" + poolExpiryDate)
	var poolSts = args[6]
	fmt.Println("the PoolSts is" + poolSts)
	var approvalDate = args[7]
	fmt.Println("the ApprovalDate is" + approvalDate)
	var approverID = args[8]
	fmt.Println("the ApproverID is" + approverID)
	var numAssets = args[9]
	fmt.Println("the NumAssets is" + numAssets)
	var poolHash = args[10]
	fmt.Println("the PoolHash is" + poolHash)
	// var loanHash = args[11]
	// fmt.Println("the loanHash is" + loanHash)
	var poolCreatedDate = args[11]
	fmt.Println("the PoolCreatedDate is" + poolCreatedDate)
	var warehouseLenderName = args[12]
	fmt.Println("the WarehouseLenderName is " + warehouseLenderName)

	createPoolStruct := CreatePool{
		PoolID:          poolID,
		PoolName:        poolName,
		PoolDesc:        poolDesc,
		Poolowner:       poolOwner,
		PoolStartedDate: poolStartedDate,
		PoolExpiryDate:  poolExpiryDate,
		PoolStatus:      poolSts,
		ApprovalDate:    approvalDate,
		ApproverID:      approverID,
		NoOfAssets:      numAssets,
		Status:          "active",
		PoolHash:        poolHash,
		// LoanHash:        loanHash,
		PoolCreatedDate:     poolCreatedDate,
		WarehouseLenderName: warehouseLenderName,
	}

	fmt.Println("the struct values are PoolID  " + createPoolStruct.PoolID + "Poolowner" + createPoolStruct.Poolowner + "PoolStartedDate" + createPoolStruct.PoolStartedDate + "ApproverID" + createPoolStruct.ApproverID)

	poolStructBytes, err := json.Marshal(createPoolStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	poolerr := stub.PutState(poolID, []byte(poolStructBytes))
	if poolerr != nil {
		fmt.Println("Could not save poolCharacterestic data to ledger", poolerr)
		return "", fmt.Errorf("Could not save poolCharacterestic data to ledger")
	}

	indexName := "createDateIndex"
	createDateKey, err := stub.CreateCompositeKey(indexName, []string{createPoolStruct.PoolCreatedDate, createPoolStruct.PoolID})
	if err != nil {
		fmt.Println("Could not index composite key for createDate", err)
		return "", fmt.Errorf("Could not index composite key for createDate map")
	}

	value := []byte{0x00}
	stub.PutState(createDateKey, value)

	indexName1 := "LenderNameIndex"
	lenderNameKey, err := stub.CreateCompositeKey(indexName1, []string{createPoolStruct.WarehouseLenderName, createPoolStruct.PoolID})
	if err != nil {
		fmt.Println("Could not index composite key for lendername", err)
		return "", fmt.Errorf("Could not index composite key for lendername map")
	}

	value1 := []byte{0x00}
	stub.PutState(lenderNameKey, value1)

	var poolEvent = "{eventType: 'Createpool', description:" + createPoolStruct.PoolID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(poolEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved createpool")
	return args[0], nil

}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(CreatePool))
	{
		fmt.Printf("Error starting CreatePool chaincode: %s", err)
	}
}
