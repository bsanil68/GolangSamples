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

type Loans struct {
	LoanID                                   string `json:"LoanID"`
	LoanContractNumber                       string `json:"LoanContractNumber"`
	City                                     string `json:"City"`
	DateOfLoanAgreement                      string `json:"DateOfLoanAgreement"`
	Lender                                   string `json:"Lender"`
	DateOfLoanApplication                    string `json:"DateOfLoanApplication"`
	SignedAndDeliveredBy                     string `json:"SignedAndDeliveredBy"`
	TypeOfLoan                               string `json:"TypeOfLoan"`
	LoanPurpose                              string `json:"LoanPurpose"`
	LoanOrFacilityAmount                     string `json:"LoanOrFacilityAmount"`
	LoanOrFacilityTermInMonths               string `json:"LoanOrFacilityTermInMonths"`
	InterestType                             string `json:"InterestType"`
	InterestChargeablePerAnnum               string `json:"InterestChargeablePerAnnum"`
	DefaultInterestRatePerAnnum              string `json:"DefaultInterestRatePerAnnum"`
	ModeOfCommunicationForInterestRateChange string `json:"ModeOfCommunicationForInterestRateChange"`
	ApplicationProcessingFee                 string `json:"ApplicationProcessingFee"`
	OtherConditions                          string `json:"OtherConditions"`
	EmiPayable                               string `json:"EmiPayable"`
	LastEMIPayable                           string `json:"LastEMIPayable"`
	DateOfCommencementOfEMI                  string `json:"DateOfCommencementOfEMI"`
	ModeOfRepayment                          string `json:"ModeOfRepayment"`
	InsurancePremium                         string `json:"InsurancePremium"`
	CalculatedLTV                            string `json:"CalculatedLTV"`
	State                                    string `json:"State"`
	OverdueAgeing                            string `json:"OverdueAgeing"`
	DefaultRating                            string `json:"DefaultRating"`
	Status                                   string `json:"Status"`
	Currency                                 string `json:"Currency"`
	LoanHash                                 string `json:"LoanHash"`
	LoanDocumentHash                         string `json:"LoanDocumentHash"`
	WarehouseLenderName                      string `json:"WarehouseLenderName"`
}

// we   dod not init  any data   when  chaincode  is called
func (t *Loans) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// invoke is called   when the client  pushes  a  post data
// to the rest  end point
func (t *Loans) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if fn == "GetLoanDetailsByLoanID" {

		result, err = GetLoanDetailsByLoanID(stub, args)
	}
	if fn == "GetLoansByLoanHash" {

		result, err = GetLoansByLoanHash(stub, args)
	}
	if fn == "GetLoansByLenderName" {

		result, err = GetLoansByLenderName(stub, args)
	}
	if fn == "GetLoandetailsforArrayofLoanHashes" {

		result, err = GetLoandetailsforArrayofLoanHashes(stub, args)
	}
	if fn == "CreateLoans" {

		result, err = CreateLoans(stub, args)
	}
	if fn == "GetAllLoans" {

		result, err = GetAllLoans(stub)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))

}

func GetAllLoans(stub shim.ChaincodeStubInterface) (string, error) {
	fmt.Println("Entering GetAllLoans")

	GetAllLoansItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all Pools")
	}
	defer GetAllLoansItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetAllLoansItr.HasNext() {
		queryResponseValue, err := GetAllLoansItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Loan Data")
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

func GetLoanDetailsByLoanID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoanDetailsByLoanID")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing AssetID")
	}

	var LoanID = args[0]
	value, err := stub.GetState(LoanID)
	if err != nil {
		fmt.Println("Could not get LoanID  with id "+LoanID+" from ledger", err)
		return "", errors.New("Missing LoanID")
	}
	return string(value), nil
}

func GetLoansByLoanHash(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoansByLoanHash")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing document")
	}

	loanHash := args[0]

	GetLoansForHashItr, err := stub.GetStateByPartialCompositeKey("hashLoanIndex", []string{loanHash})
	if err != nil {
		return "", fmt.Errorf("Could not get all Loans for Hash")
	}
	defer GetLoansForHashItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"LoansForPool\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetLoansForHashItr.HasNext() {
		queryResponse, err := GetLoansForHashItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Hash Loans Data")
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
		returnedLoanID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s hashID:%s loanID:%s\n", objectType, returnedHashID, returnedLoanID)

		value, err := stub.GetState(returnedLoanID)
		if err != nil {
			fmt.Println("Could not get details for id "+returnedLoanID+" from ledger", err)
			return "", errors.New("Missing LoanID")
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

func GetLoansByLenderName(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoansByLenderName")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing document")
	}

	lenderName := args[0]

	GetLoansForLenderNameItr, err := stub.GetStateByPartialCompositeKey("lenderNameIndex", []string{lenderName})
	if err != nil {
		return "", fmt.Errorf("Could not get all Loans for this lender")
	}
	defer GetLoansForLenderNameItr.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("{\"LoansForPool\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for GetLoansForLenderNameItr.HasNext() {
		queryResponse, err := GetLoansForLenderNameItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Loans Data")
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
		returnedLenderName := compositeKeyParts[0]
		returnedLoanID := compositeKeyParts[1]

		fmt.Printf("- found a  from index:%s LenderName:%s loanID:%s\n", objectType, returnedLenderName, returnedLoanID)

		value, err := stub.GetState(returnedLoanID)
		if err != nil {
			fmt.Println("Could not get details for id "+returnedLoanID+" from ledger", err)
			return "", errors.New("Missing LoanID")
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

func GetLoandetailsforArrayofLoanHashes(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetLoandetailsforArrayofLoanHashes")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing array")
	}

	str := args[0]
	var strArray = strings.Split(str, "#")
	var buffer bytes.Buffer

	buffer.WriteString("{\"LoansForPool\":")

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	if len(strArray) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing document")
	}
	for i := 0; i < len(strArray); i++ {

		loanHash := strArray[i]

		GetLoansForHashItr, err := stub.GetStateByPartialCompositeKey("hashLoanIndex", []string{loanHash})
		if err != nil {
			return "", fmt.Errorf("Could not get all Loans for Hash")
		}
		defer GetLoansForHashItr.Close()

		// buffer is a JSON array containing QueryResults

		for GetLoansForHashItr.HasNext() {
			queryResponse, err := GetLoansForHashItr.Next()
			if err != nil {
				return "", fmt.Errorf("Could not get next Hash Loans Data")
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
			returnedLoanID := compositeKeyParts[1]

			fmt.Printf("- found a  from index:%s hashID:%s loanID:%s\n", objectType, returnedHashID, returnedLoanID)

			value, err := stub.GetState(returnedLoanID)
			if err != nil {
				fmt.Println("Could not get details for id "+returnedLoanID+" from ledger", err)
				return "", errors.New("Missing LoanID")
			}
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(value))
			bArrayMemberAlreadyWritten = true
		}

		//fmt.Printf("- GetAllPoolAssigneeMaps queryResult:\n%s\n", buffer.String())
	}

	buffer.WriteString("]")
	buffer.WriteString("}")

	var basket string
	basket = buffer.String()
	return string(basket), nil
}

func CreateLoans(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CreateLoans")

	if len(args) < 1 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var loan_id = args[0]
	fmt.Println("the loanId is" + loan_id)
	var loan_contract_no = args[1]
	fmt.Println("the loan_contract_no is" + loan_contract_no)
	var city = args[2]
	fmt.Println("the city is" + city)
	var date_of_loan_agr = args[3]
	fmt.Println("the date_of_loan_agr is" + date_of_loan_agr)
	var leander = args[4]
	fmt.Println("the customerid is" + leander)
	var date_of_laon_application = args[5]
	fmt.Println("the date_of_laon_application is" + date_of_laon_application)
	var signed_and_deliveredBy = args[6]
	fmt.Println("the signed_and_deliveredBy is" + signed_and_deliveredBy)
	var type_loan = args[7]
	fmt.Println("the type_loan is" + type_loan)
	var loan_purpose = args[8]
	fmt.Println("the loan_purpose is" + loan_purpose)
	var loan_amount = args[9]
	fmt.Println("the loan_amount is" + loan_amount)
	var loan_facility_months = args[10]
	fmt.Println("the loan_facility_months is" + loan_facility_months)
	var interest_type = args[11]
	fmt.Println("the interest_type is" + interest_type)
	var intrest_per_anum = args[12]
	fmt.Println("the intrest_per_anum is" + intrest_per_anum)
	var default_intrest_per_anum = args[13]
	fmt.Println("the default_intrest_per_anum is" + default_intrest_per_anum)
	var mode_cmm_int_rate_charge = args[14]
	fmt.Println("the mode_cmm_int_rate_charge is" + mode_cmm_int_rate_charge)
	var application_process_fee = args[15]
	fmt.Println("the application_process_fee is" + application_process_fee)
	var other_conditions = args[16]
	fmt.Println("the other_conditions is" + other_conditions)
	var emi_payble = args[17]
	fmt.Println("the emi_payble is" + emi_payble)
	var last_emi_payble = args[18]
	fmt.Println("the last_emi_payble is" + last_emi_payble)
	var date_cmm_emi = args[19]
	fmt.Println("the date_cmm_emi is" + date_cmm_emi)
	var mode_repayment = args[20]
	fmt.Println("the mode_repayment is" + mode_repayment)
	var insurance_prem = args[21]
	fmt.Println("the insurance_prem is" + insurance_prem)
	var cal_ltv = args[22]
	fmt.Println("the cal_ltv is" + cal_ltv)
	var state = args[23]
	fmt.Println("the state is" + state)
	var overDue_ageing = args[24]
	fmt.Println("the overDue_ageing is" + overDue_ageing)
	var default_rating = args[25]
	fmt.Println("the default_rating is" + default_rating)
	var currency = args[26]
	fmt.Println("the currency is" + currency)
	var loanHash = args[27]
	fmt.Println("the loanHash is" + loanHash)
	var loanDocumentHash = args[28]
	fmt.Println("the loan document hash is " + loanDocumentHash)
	var warehouseLenderName = args[29]
	fmt.Println("the warehouse lender name is " + warehouseLenderName)

	//assigning to struct the variables
	LoansStruct := Loans{
		LoanID:                                   loan_id,
		LoanContractNumber:                       loan_contract_no,
		City:                                     city,
		DateOfLoanAgreement:                      date_of_loan_agr,
		Lender:                                   leander,
		DateOfLoanApplication:                    date_of_laon_application,
		SignedAndDeliveredBy:                     signed_and_deliveredBy,
		TypeOfLoan:                               type_loan,
		LoanPurpose:                              loan_purpose,
		LoanOrFacilityAmount:                     loan_amount,
		LoanOrFacilityTermInMonths:               loan_facility_months,
		InterestType:                             interest_type,
		InterestChargeablePerAnnum:               intrest_per_anum,
		DefaultInterestRatePerAnnum:              default_intrest_per_anum,
		ModeOfCommunicationForInterestRateChange: mode_cmm_int_rate_charge,
		ApplicationProcessingFee:                 application_process_fee,
		OtherConditions:                          other_conditions,
		EmiPayable:                               emi_payble,
		LastEMIPayable:                           last_emi_payble,
		DateOfCommencementOfEMI:                  date_cmm_emi,
		ModeOfRepayment:                          mode_repayment,
		InsurancePremium:                         insurance_prem,
		CalculatedLTV:                            cal_ltv,
		State:                                    state,
		OverdueAgeing:                            overDue_ageing,
		DefaultRating:                            default_rating,
		Status:                                   "active",
		Currency:                                 currency,
		LoanHash:                                 loanHash,
		LoanDocumentHash:                         loanDocumentHash,
		WarehouseLenderName:                      warehouseLenderName,
	}

	fmt.Println("the struct values are loanId " + LoansStruct.LoanID + "loanContractNumber" + LoansStruct.LoanContractNumber + "city" + LoansStruct.City)

	loanStructBytes, err := json.Marshal(LoansStruct)
	if err != nil {
		fmt.Println("Could not masrhall data from struct", err)
		return "", fmt.Errorf("Could marshal data from struct")

	}
	loanerr := stub.PutState(loan_id, []byte(loanStructBytes))
	if loanerr != nil {
		fmt.Println("Could not save loans data to ledger", loanerr)
		return "", fmt.Errorf("Could not save loans data to ledger")
	}
	//creating Composite Key index to be able to query based on it
	indexName := "hashLoanIndex"
	hashLoanKey, err := stub.CreateCompositeKey(indexName, []string{LoansStruct.LoanHash, LoansStruct.LoanID})
	if err != nil {
		fmt.Println("Could not index composite key for pool and loan", err)
		return "", fmt.Errorf("Could not index composite key for pool loan")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(hashLoanKey, value)

	indexName1 := "lenderNameIndex"
	lenderNameKey, err := stub.CreateCompositeKey(indexName1, []string{LoansStruct.WarehouseLenderName, LoansStruct.LoanID})
	if err != nil {
		fmt.Println("Could not index composite key for pool and loan", err)
		return "", fmt.Errorf("Could not index composite key for pool loan")
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the collection.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value1 := []byte{0x00}
	stub.PutState(lenderNameKey, value1)

	var loanEvent = "{eventType: 'CreateLoans', description:" + LoansStruct.LoanID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(loanEvent))
	if err != nil {
		return "", fmt.Errorf("Could not set event hub")
	}
	fmt.Println("Successfully saved loans data")
	return args[0], nil

}

// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(Loans))
	{
		fmt.Printf("Error starting enrollcustomer chaincode: %s", err)
	}
}
