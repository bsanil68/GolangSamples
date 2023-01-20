package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type BankAccountDetailsCC struct {
	UserId              string `json:"UserId"`
	AccountNumber       string `json:"AccountNumber"`
	BankName            string `json:"BankName"`
	BankIFSCOrIBankCode string `json:"BankIFSCOrIBankCode"`
	Currency            string `json:"Currency"`
}

func (t *BankAccountDetailsCC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *BankAccountDetailsCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetBankAccountDetails" {

		result, err = GetBankAccountDetails(stub, args)
	}
	if fn == "SaveBankAccountDetails" {

		result, err = SaveBankAccountDetails(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetBankAccountDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetBankAccountDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing userId")
	}

	var userID = args[0]
	value, err := stub.GetState(userID)
	if err != nil {
		fmt.Println("Couldn't get user_ID  with id "+userID+" from ledger", err)
		return "", errors.New("Missing userID")
	}
	return string(value), nil
}

func SaveBankAccountDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering SaveBankAccountDetails")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var userID = args[0]
	fmt.Println("the userId is" + userID)
	var accountNumber = args[1]
	fmt.Println("the account number is" + accountNumber)
	var bankName = args[2]
	fmt.Println("the bank name is" + bankName)
	var bankIFSCorIBankCode = args[3]
	fmt.Println("the bank IFSC or IBank code is" + bankIFSCorIBankCode)
	var userCurrency = args[4]
	fmt.Println("the currency is" + userCurrency)

	BankAccountDetailsCCStruct := BankAccountDetailsCC{
		UserId:              userID,
		AccountNumber:       accountNumber,
		BankName:            bankName,
		BankIFSCOrIBankCode: bankIFSCorIBankCode,
		Currency:            userCurrency,
	}

	fmt.Println("the struct values are userId  " + BankAccountDetailsCCStruct.UserId + "accountNumber" + BankAccountDetailsCCStruct.AccountNumber + "bankName" + BankAccountDetailsCCStruct.BankName + "bankIFSCOrIBankCode" + BankAccountDetailsCCStruct.BankIFSCOrIBankCode + "currency" + BankAccountDetailsCCStruct.Currency)

	bankAccountStructBytes, err := json.Marshal(BankAccountDetailsCCStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	bankAccounterr := stub.PutState(userID, []byte(bankAccountStructBytes))
	if bankAccounterr != nil {
		fmt.Println("Couldn't save bankAccountCharacterestic data to ledger", bankAccounterr)
		return "", fmt.Errorf("Couldn't save bankAccountCharacterestic data to ledger")
	}

	var bankAccountEvent = "{eventType: 'BankAccountDetailsCC', description:" + BankAccountDetailsCCStruct.UserId + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(bankAccountEvent))
	if err != nil {
		return "", fmt.Errorf("Couldn't set event hub")
		fmt.Println("Successfully saved BankAccountDetailsCC")
	}
	return args[0], nil
}

func main() {
	err := shim.Start(new(BankAccountDetailsCC))
	{
		fmt.Printf("Error starting BankAccountDetailsCC chaincode: %s", err)
	}
}
