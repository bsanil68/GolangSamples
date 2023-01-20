export CORE_PEER_LOCALMSPID="Org0MSP"
 export CORE_PEER_TLS_ROOTCERT_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/peers/peer0.org0/tls/ca.crt
 export CORE_PEER_MSPCONFIGPATH=/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/users/Admin@org0/msp
 export CORE_PEER_ADDRESS=peer0:7051
export  CORE_PEER_TLS_CLIENTROOTCAS_FILES=/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/users/Admin@org0/tls/ca.crt
export  CORE_PEER_TLS_CLIENTCERT_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/users/Admin@org0/tls/client.crt
export CORE_PEER_TLS_CLIENTKEY_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/users/Admin@org0/tls/client.key
export CORE_PEER_TLS_ENABLED=true

export ORDERER_GENERAL_TLS_PRIVATEKEY=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/tls/server.key
export ORDERER_GENERAL_TLS_CERTIFICATE=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/tls/server.crt
export ORDERER_GENERAL_TLS_ROOTCAS=[/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/tls/ca.crt,/root/bcnetwork/conf/crypto-config/peerOrganizations/org0/peers/peer0.org0/tls/ca.crt]

export $CORE_PEER_TLS_ENABLED=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/tls/ca.crt

export Orderer_CA=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem

echo "LoanID is:" $1
LoanID=$1

echo "LoanContractNumber is:" $2
LoanContractNumber=$2

echo "City is:" $3
City=$3

echo "DateOfLoanAgreement is:" $4
DateOfLoanAgreement=$4

echo "Lender is:" $5
Lender=$5

echo "DateOfLoanApplication is:" $6
DateOfLoanApplication=$6

echo "SignedAndDeliveredBy is:" $7
SignedAndDeliveredBy=$7

echo "TypeOfLoan is:" $8
TypeOfLoan=$8

echo "LoanPurpose is:" ${9}
LoanPurpose=${9}

echo "LoanOrFacilityAmount is:" ${10}
LoanOrFacilityAmount=${10}

echo "LoanOrFacilityTermInMonths is:" ${11}
LoanOrFacilityTermInMonths=${11}

echo "InterestType is:" ${12}
InterestType=${12}

echo "InterestChargeablePerAnnum is:" ${13}
InterestChargeablePerAnnum=${13}

echo "DefaultInterestRatePerAnnum is:" ${14}
DefaultInterestRatePerAnnum=${14}

echo "ModeOfCommunicationForInterestRateChange is:" ${15}
ModeOfCommunicationForInterestRateChange=${15}

echo "ApplicationProcessingFee is:" ${16}
ApplicationProcessingFee=${16}

echo "OtherConditions is:" ${17}
OtherConditions=${17}

echo "EmiPayable is:" ${18}
EmiPayable=${18}

echo "LastEMIPayable is:" ${19}
LastEMIPayable=${19}

echo "DateOfCommencementOfEMI is:" ${20}
DateOfCommencementOfEMI=${20}

echo "ModeOfRepayment is:" ${21}
ModeOfRepayment=${21}

echo "InsurancePremium is:" ${22}
InsurancePremium=${22}

echo "CalculatedLTV is:" ${23}
CalculatedLTV=${23}

echo "State is:" ${24}
State=${24}

echo "OverdueAgeing is:" ${25}
OverdueAgeing=${25}

echo "DefaultRating is:" ${26}
DefaultRating=${26}

echo "Currency:" ${27}
Currency=${27}

echo "LoanHash:" ${28}
LoanHash=${28}

peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n Loans3CC -c '{"Args":["CreateLoans","'"$LoanID"'","'"$LoanContractNumber"'","'"$City"'","'"$DateOfLoanAgreement"'","'"$Lender"'","'"$DateOfLoanApplication"'","'"$SignedAndDeliveredBy"'","'"$TypeOfLoan"'","'"$LoanPurpose"'","'"$LoanOrFacilityAmount"'","'"$LoanOrFacilityTermInMonths"'","'"$InterestType"'",
"'"$InterestChargeablePerAnnum"'","'"$DefaultInterestRatePerAnnum"'","'"$ModeOfCommunicationForInterestRateChange"'","'"$ApplicationProcessingFee"'","'"$OtherConditions"'","'"$EmiPayable"'","'"$LastEMIPayable"'","'"$DateOfCommencementOfEMI"'","'"$ModeOfRepayment"'","'"$InsurancePremium"'","'"$CalculatedLTV"'","'"$State"'","'"$OverdueAgeing"'","'"$DefaultRating"'","'"$Currency"'","'"$LoanHash"'"]}'
