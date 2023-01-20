package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Document struct {
	DocID     string `json:"DocID"`
	DocName   string `json:"DocName"`
	DocType   string `json:"DocType"`
	DocTag    string `json:"DocTag"`
	DocPath   string `json:"DocPath"`
	HashValue string `json:"HashValue"`
	PoolID    string `json:"PoolID"`
	OwnerID   string `json:"OwnerID"`
	OwnerTag  string `json:"OwnerTag"`
}

func (t *Document) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/***

invokes  the   respective chaincode

**/
func (t *Document) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "SaveDocumentDetails" {
		result, err = SaveDocumentDetails(stub, args)
	} else if fn == "GetDocumentByID" {
		result, err = GetDocumentByID(stub, args)
	} else if fn == "GetAllDocs" {
		result, err = GetAllDocs(stub)
	} else if fn == "GetDocumentDetailsByPoolIDAndDocTag" {
		result, err = GetDocumentDetailsByPoolIDAndDocTag(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}

/*func createKey(args []string) (string, error) {

	sourceData := []string{args[0], args[1]}
	var result = strings.Join(sourceData, "#")
	return result, nil
}*/

func GetDocumentByID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	var docID = args[0]

	value, err := stub.GetState(docID)
	if err != nil {
		return "", fmt.Errorf("Failed to get document: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("document not found: %s", args[0])
	}
	return string(value), nil
}

func SaveDocumentDetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering saving document data ")
	/*
		var result string
		var validatedValue string
	*/

	if len(args) < 9 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	/** validate  the  loanID and pool  id  does  not  exist  in ledger false

	validatedValue, err := validateDocAndTypeIsMapped(stub, args)
	if err != nil {
		fmt.Println("error in calling  validation", err)
		return "", fmt.Errorf("validation functions not called")
	}

	if validatedValue == "notassigned" {
	*/
	// json args input musts match  the case and  spelling exactly
	// get all the  arguments
	var docID = args[0]
	fmt.Println("the docID is" + docID)
	var docName = args[1]
	fmt.Println("the docName is" + docName)
	var docType = args[2]
	fmt.Println("the docType is" + docType)
	var docTag = args[3]
	fmt.Println("the docTag is" + docTag)
	var docPath = args[4]
	fmt.Println("the docGrp is" + docPath)
	var hashValue = args[5]
	fmt.Println("the HashValue is" + hashValue)
	var poolID = args[6]
	fmt.Println("the pool ID is" + poolID)
	var ownerID = args[7]
	fmt.Println("the owner ID is" + ownerID)
	var ownerTag = args[8]
	fmt.Println("the ownerTag is" + ownerTag)

	//assigning to struct the variables
	DocStruct := Document{
		DocID:     docID,
		DocName:   docName,
		DocType:   docType,
		DocTag:    docTag,
		DocPath:   docPath,
		HashValue: hashValue,
		PoolID:    poolID,
		OwnerID:   ownerID,
		OwnerTag:  ownerTag,
	}

	docStructBytes, err := json.Marshal(DocStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	doctypeerr := stub.PutState(docID, []byte(docStructBytes))
	if doctypeerr != nil {
		fmt.Println("Could not save  data to ledger", doctypeerr)
		return "", fmt.Errorf("Could not save  data to ledger")
	}
	/*
		keyComb := []string{args[6], args[3]}
		keyValue, err := createKey(keyComb)
		if err != nil {
			return "", fmt.Errorf("Could not create key combination poolId and docTag")
		}*/

	indexName := "doctypeIndex"
	docKey, err := stub.CreateCompositeKey(indexName, []string{DocStruct.PoolID, DocStruct.DocTag, DocStruct.DocID})
	if err != nil {
		fmt.Println("Could not index composite key for document and type", err)
		return "", fmt.Errorf("Could not index composite key for doc type map")
	}

	value := []byte{0x00}
	stub.PutState(docKey, value)

	var customEvent = "{eventType: 'CreateDocTypeMap', description: docID: " + DocStruct.DocID + " docType: " + DocStruct.DocType + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	/*
			result = "Validated"
		}
		if validatedValue == "assigned" {
			result = "notvalidated"

		}
		if validatedValue == "crazyerror" {
			result = "notvalidated"

		}
		if validatedValue == "spliterror" {
			result = "notvalidated"

		} */

	fmt.Println("Successfully saved doc and doctype map")
	return args[0], nil

}

func GetDocumentDetailsByPoolIDAndDocTag(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetDocumentDetailsByPoolIDandDcoTag")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments: PoolId and DocTag are required")
		return "", fmt.Errorf("Missing arguments")
	}

	var poolID = args[0]
	var docTag = args[1]

	/*keysInfo := []string{args[0], args[1]}
	keyValue, err := createKey(keysInfo)
	if err != nil {
		return "", fmt.Errorf("Could not create key combination from poolId and DocTag")
	}*/
	poolDocumentItr, err := stub.GetStateByPartialCompositeKey("doctypeIndex", []string{poolID, docTag})
	if err != nil {
		return "", fmt.Errorf("Could not get documents for docId and docType")
	}
	defer poolDocumentItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"DocumentDetails\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for poolDocumentItr.HasNext() {
		queryResponse, err := poolDocumentItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next document data")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite Key Parts")
		}
		returnedPoolID := compositeKeyParts[0]
		returnedDocTag := compositeKeyParts[1]
		returnedDocID := compositeKeyParts[2]

		fmt.Printf("- found from index:%s poolIdDocTypeCombination:%s DocId:%s\n", objectType, returnedPoolID, returnedDocTag, returnedDocID)

		/*s := strings.Split(returnedPoolIDDocTagComb, "#")
		splitPoolID, splitDocTag := s[0], s[1]
		fmt.Printf("Split poolID:%S, docTag:%S", splitPoolID, splitDocTag)*/

		value, err := stub.GetState(returnedDocID)
		if err != nil {
			fmt.Println("Couldn't get document for DocID"+returnedDocID+"from ledger", err)
			return "", errors.New("Missing DocID")
		}
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

//assign  document and   document  type  as  composite  key  and save the  others   based on this
/**
func AddDocAndDocType(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering saving document data ")

	var result string
	var validatedValue string
	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	//  validate  the  loanID and pool  id  does  not  exist  in ledger false

	validatedValue, err := validateDocAndTypeIsMapped(stub, args)
	if err != nil {
		fmt.Println("error in calling  validation", err)
		return "", fmt.Errorf("validation functions not called")
	}

	if validatedValue == "notassigned" {

		// json args input musts match  the case and  spelling exactly
		// get all the  arguments

		var docID = args[0]
		fmt.Println("the docID is" + docID)
		var docTypeID = args[1]
		fmt.Println("the docTypeID is" + docTypeID)
		var docName = args[2]
		fmt.Println("the docName is" + docName)
		var docTag = args[3]
		fmt.Println("the docTag is" + docTag)
		var ownerid = args[4]
		fmt.Println("the ownerid is" + ownerid)
		var statusTag = args[5]
		fmt.Println("the statusTag is" + statusTag)
		var docHash = args[6]
		fmt.Println("the docHash is" + docHash)
		var docGrp = args[7]
		fmt.Println("the docGrp is" + docGrp)
		var ownerTag = args[8]
		fmt.Println("the ownerTag is" + ownerTag)

		//assigning to struct the variables
		DocStruct := Document{
			DocID:        docID,
			DocType:   docTypeID,
			DocName:   docName,
			DocTag:    docTag,
			OwnerID:   ownerid,
			StatusTag: statusTag,
			HashValue: docHash,
			DocGroup:  docGrp,
			OwnerTag:  ownerTag,
		}

		docStructBytes, err := json.Marshal(DocStruct)

		if err != nil {
			fmt.Println("Could not masrhall data from struct", err)
			return "", fmt.Errorf("Could marshal data from struct")

		}
		doctypeerr := stub.PutState(docID, []byte(docStructBytes))
		if doctypeerr != nil {
			fmt.Println("Could not save  data to ledger", doctypeerr)
			return "", fmt.Errorf("Could not save  data to ledger")
		}

		//creating Composite Key index to be able to query based on it
		indexName := "doctypeIndex"
		docKey, err := stub.CreateCompositeKey(indexName, []string{DocStruct.ID, DocStruct.DocType})
		if err != nil {
			fmt.Println("Could not index composite key for documnet and type", err)
			return "", fmt.Errorf("Could not index composite key for doc type map")
		}
		//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		stub.PutState(docKey, value)

		var customEvent = "{eventType: 'CreateDocTypeMap', description: docID: " + DocStruct.ID + " docTypeID: " + DocStruct.DocType + "' Successfully created'}"
		err = stub.SetEvent("evtSender", []byte(customEvent))
		if err != nil {
			return "", fmt.Errorf("Could not set event hub")
		}
		result = "Validated"
	}
	if validatedValue == "assigned" {
		result = "notvalidated"

	}
	if validatedValue == "crazyerror" {
		result = "notvalidated"

	}
	if validatedValue == "spliterror" {
		result = "notvalidated"

	}
	fmt.Println("Successfully saved doc and doctype map")
	return result, nil

}
**/
/*
func validateDocAndTypeIsMapped(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering doc type map check   ")

	var res string
	// docid
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing document")
	}

	docID := args[0]

	fmt.Println("- check  if  document  is already mapped   ", docID)

	// Query the doc~doctypeIndex index by docid
	// This will execute a key range query on all keys starting with ''
	docIterator, err := stub.GetStateByPartialCompositeKey("doctypeIndex", []string{docID})
	if err != nil {
		res = "notassigned"
		return res, fmt.Errorf("did not get  any document ,Validated ")
	}
	defer docIterator.Close()

	// Iterate through result set and check
	var i int
	for i = 0; docIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the loan id from the composite key
		responseRange, err := docIterator.Next()
		if err != nil {
			res = "crazyerror"
			return res, fmt.Errorf("doc iterator, doc found  but doctype not  found crazy issues")

		}

		// get the loanid and name from loanindex composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			fmt.Println("- split error   ")
			res = "spliterror"
			return res, fmt.Errorf("splitting error ")
		}
		returnedDocID := compositeKeyParts[0]
		returnedDocTypeID := compositeKeyParts[1]
		fmt.Printf("- found a  from index:%s loanid:%s poolid:%s\n", objectType, returnedDocID, returnedDocTypeID)
		res = "assigned"
	}

	return res, nil
}
*/
func GetAllDocs(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllpools")

	GetAllDocsItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Pools")
	}
	defer GetAllDocsItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"AllDocumentDetails\":")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllDocsItr.HasNext() {
		queryResponseValue, err := GetAllDocsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Docs Data")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponseValue.GetValue()))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}

func main() {
	if err := shim.Start(new(Document)); err != nil {
		fmt.Printf("Error starting Document chaincode: %s", err)
	}
}
