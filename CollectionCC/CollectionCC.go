package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Collection struct {
	CollectionID                  string `json:"CollectionID"`
	LoanID                        string `json:"LoanID"`
	InstallmentNumber             string `json:"InstallmentNumber"`
	EmiPaid                       string `json:"EmiPaid"`
	InterestAmountRepaid          string `json:"InterestAmountRepaid"`
	PrincipalAmountRepaid         string `json:"PrincipalAmountRepaid"`
	OutstandingPrincipalBalance   string `json:OutstandingPrincipalBalance`
	OverdueEMINumbers1To30Days    string `json:OverdueEMINumbers1To30Days`
	OverdueEMINumbers31To60Days   string `json:OverdueEMINumbers31To60Days`
	OverdueEMINumbers61To90Days   string `json:OverdueEMINumbers61To90Days`
	OverdueEMINumbers91To120Days  string `json:OverdueEMINumbers91To120Days`
	OverdueEMINumbers121To180Days string `json:OverdueEMINumbers121To180Days`
	OverdueEMINumbers180PlusDays  string `json:OverdueEMINumbers180PlusDays`
	Month                         string `json:"Month"`
	Currency                      string `json:"Currency"`
	CollectionsHash               string `json:"CollectionsHash"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *Collection) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *Collection) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "SaveCollection" {

		result, err = SaveCollection(stub, args)
	}
	if fn == "GetAllCollectionsForHash" {

		result, err = GetAllCollectionsForHash(stub, args)
	}
	if fn == "GetAllCollectionsForLoan" {

		result, err = GetAllCollectionsForLoan(stub, args)
	}
	if fn == "GetCollectionsforArrayofLoans" {

		result, err = GetCollectionsforArrayofLoans(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

//Get Collections for a Pool Function
func GetAllCollectionsForHash(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetAllCollectionsForHash")
	//var emptyBuf bytes.Buffer

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var Collections_Hash = args[0]

	hashCollectionsItr, err := stub.GetStateByPartialCompositeKey("hashCollectionsIndex", []string{Collections_Hash})
	if err != nil {
		return "", fmt.Errorf("Could not get Collections for Hash")
	}
	defer hashCollectionsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"CollectionsForHash\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for hashCollectionsItr.HasNext() {
		queryResponse, err := hashCollectionsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next hash Collections Data")
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
		returnedHashID := compositeKeyParts[0]
		returnedCollectionID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s hashID:%s collectionID:%s\n", objectType, returnedHashID, returnedCollectionID)

		value, err := stub.GetState(returnedCollectionID)
		if err != nil {
			fmt.Println("Could not get details for id "+returnedCollectionID+" from ledger", err)
			return "", errors.New("Missing CollectionID")
		}
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetAllCollectionsForLoan(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetAllCollectionsForLoan")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	var Loan_ID = args[0]
	loanCollectionsItr, err := stub.GetStateByPartialCompositeKey("loanCollectionsIndex", []string{Loan_ID})
	if err != nil {
		return "", fmt.Errorf("Could not get Collections for Loan")
	}
	defer loanCollectionsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"CollectionsForLoan\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for loanCollectionsItr.HasNext() {
		queryResponse, err := loanCollectionsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Pool Loans Data")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		returnedLoanID := compositeKeyParts[0]
		returnedCollectionID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s loanid:%s collectionID:%s\n", objectType, returnedLoanID, returnedCollectionID)

		value, err := stub.GetState(returnedCollectionID)
		if err != nil {
			fmt.Println("Could not get details for id "+returnedCollectionID+" from ledger", err)
			return "", errors.New("Missing CollectionID")
		}
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func GetCollectionsforArrayofLoans(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetCollectionsforArrayofLoans")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing array")
	}

	str := args[0]
	var strArray = strings.Split(str, "#")
	var buffer bytes.Buffer

	if len(strArray) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", fmt.Errorf("Invalid number of arguments")
	}

	buffer.WriteString("{\"CollectionsForPool\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for i := 0; i < len(strArray); i++ {

		loan_ID := strArray[i]

		loanCollectionsItr, err := stub.GetStateByPartialCompositeKey("loanCollectionsIndex", []string{loan_ID})
		if err != nil {
			return "", fmt.Errorf("Could not get Collections for Loan")
		}
		defer loanCollectionsItr.Close()

		// buffer is a JSON array containing QueryResults
		
		for loanCollectionsItr.HasNext() {
			queryResponse, err := loanCollectionsItr.Next()
			if err != nil {
				return "", fmt.Errorf("Could not get next Pool Loans Data")
			}
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}

			objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
			if err != nil {
				return "", fmt.Errorf("Could not get Composite Key Parts")
			}
			returnedLoanID := compositeKeyParts[0]
			returnedCollectionID := compositeKeyParts[1]

			fmt.Printf("- found a  from index:%s loanid:%s collectionID:%s\n", objectType, returnedLoanID, returnedCollectionID)

			value, err := stub.GetState(returnedCollectionID)
			if err != nil {
				fmt.Println("Could not get details for id "+returnedCollectionID+" from ledger", err)
				return "", errors.New("Missing CollectionID")
			}
			buffer.WriteString(string(value))
			bArrayMemberAlreadyWritten = true
		}
		
	}

	buffer.WriteString("]")
	buffer.WriteString("}")

	var basket string
	basket = buffer.String()
	return string(basket), nil
}

func SaveCollection(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) < 16 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 16")
	}
	var collection_ID = args[0]
	fmt.Println("the Collection is" + collection_ID)
	var loan_ID = args[1]
	fmt.Println("the Asset is" + loan_ID)
	var installment_number = args[2]
	fmt.Println("the Installment Number is" + installment_number)
	var emi_paid = args[3]
	fmt.Println("the EMI Paid is" + emi_paid)
	var interest_amount_repaid = args[4]
	fmt.Println("the Interest Amount Repaid is" + interest_amount_repaid)
	var principal_amount_repaid = args[5]
	fmt.Println("the Principal Amount Repaid is" + principal_amount_repaid)
	var outstanding_principal_balance = args[6]
	fmt.Println("the Outstanding Principal Balance is" + outstanding_principal_balance)
	var overdue_EMI_numbers_1To30_days = args[7]
	fmt.Println("the Overdue EMI Numbers in 1 to 30 days is" + overdue_EMI_numbers_1To30_days)
	var overdue_EMI_numbers_31To60_days = args[8]
	fmt.Println("the Overdue EMI Numbers in 31 to 60 days is" + overdue_EMI_numbers_31To60_days)
	var overdue_EMI_numbers_61To90_days = args[9]
	fmt.Println("the Overdue EMI Numbers in 61 to 90 days is" + overdue_EMI_numbers_61To90_days)
	var overdue_EMI_numbers_91To120_days = args[10]
	fmt.Println("the Overdue EMI Numbers in 91 to 120 days is" + overdue_EMI_numbers_91To120_days)
	var overdue_EMI_numbers_121To180_days = args[11]
	fmt.Println("the Overdue EMI Numbers in 121 to 180 days is" + overdue_EMI_numbers_121To180_days)
	var overdue_EMI_numbers_180plus_days = args[12]
	fmt.Println("the Overdue EMI Numbers in 180 plus days is" + overdue_EMI_numbers_180plus_days)
	var month = args[13]
	fmt.Println("the collection month is" + month)
	var currency = args[14]
	fmt.Println("the currency is" + currency)
	var collectionsHash = args[15]
	fmt.Println("the collectionsHash is" + collectionsHash)

	collectionStruct := Collection{
		CollectionID:                  collection_ID,
		LoanID:                        loan_ID,
		InstallmentNumber:             installment_number,
		EmiPaid:                       emi_paid,
		InterestAmountRepaid:          interest_amount_repaid,
		PrincipalAmountRepaid:         principal_amount_repaid,
		OutstandingPrincipalBalance:   outstanding_principal_balance,
		OverdueEMINumbers1To30Days:    overdue_EMI_numbers_1To30_days,
		OverdueEMINumbers31To60Days:   overdue_EMI_numbers_31To60_days,
		OverdueEMINumbers61To90Days:   overdue_EMI_numbers_61To90_days,
		OverdueEMINumbers91To120Days:  overdue_EMI_numbers_91To120_days,
		OverdueEMINumbers121To180Days: overdue_EMI_numbers_121To180_days,
		OverdueEMINumbers180PlusDays:  overdue_EMI_numbers_180plus_days,
		Month:                         month,
		Currency:                      currency,
		CollectionsHash:               collectionsHash,
	}

	fmt.Println("the struct values are LoanID " + collectionStruct.LoanID + " Installment Number " + collectionStruct.InstallmentNumber)

	collectionStructBytes, err := json.Marshal(collectionStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	collectionerr := stub.PutState(collection_ID, []byte(collectionStructBytes))
	if collectionerr != nil {
		fmt.Println("Could not save Collection data to ledger", collectionerr)
		return "", fmt.Errorf("Could not save Collection data to ledger")
	}

	//creating Composite Key index to be able to query based on it
	indexName1 := "loanCollectionsIndex"
	collectionLoanIndexKey, err := stub.CreateCompositeKey(indexName1, []string{collectionStruct.LoanID, collectionStruct.CollectionID})
	if err != nil {
		fmt.Println("Could not index composite key for collection", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value1 := []byte{0x00}
	stub.PutState(collectionLoanIndexKey, value1)

	//creating Composite Key index to be able to query based on it
	indexName2 := "hashCollectionsIndex"
	collectionHashIndexKey, err := stub.CreateCompositeKey(indexName2, []string{collectionStruct.CollectionsHash, collectionStruct.CollectionID})
	if err != nil {
		fmt.Println("Could not index composite key for collection", err)
		return "", fmt.Errorf("Could not index composite key for collection")
	}
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value2 := []byte{0x00}
	stub.PutState(collectionHashIndexKey, value2)

	var collectionEvent = "{eventType: 'Collection', description: CollectionsHash: " + collectionStruct.CollectionsHash + " and LoanID: " + collectionStruct.LoanID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(collectionEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved collection")
	return args[0], nil

}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(Collection))
	{
		fmt.Printf("Error starting Collection chaincode: %s", err)
	}
}
