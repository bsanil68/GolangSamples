ISID=$1
ISDealId=$2
ISMaturityDate=$3
ISNodInitialAccrualPeriod=$4
ISNodAccrualPeriod=$5
ISUpdatedBy=$6
ISUpdationDate=$7
ISReserveFundBalance=$8
ISCertificateBalance=$9
ISCutoffLoanCount=${10}
ISCutoffLoanBalance=${11}
ISCutoffAverageBalance=${12}
ISCutoffWtdAvgRate=${13}
ISCutoffWtdAvgTerm=${14}

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n InitialSetupCC -c '{"Args":["SaveInitialSetup","'"$ISID"'","'"$ISDealId"'","'"$ISMaturityDate"'","'"$ISNodInitialAccrualPeriod"'","'"$ISNodAccrualPeriod"'","'"$ISUpdatedBy"'","'"$ISUpdationDate"'","'"$ISReserveFundBalance"'","'"$ISCertificateBalance"'","'"$ISCutoffLoanCount"'","'"$ISCutoffLoanBalance"'","'"$ISCutoffAverageBalance"'","'"$ISCutoffWtdAvgRate"'",
"'"$ISCutoffWtdAvgTerm"'"]}'
