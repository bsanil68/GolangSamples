IRID=$1
IRData=$2
IRMonth=$3
IRYear=$4
IRUpdatedBy=$5
IRUpdationDate=$6
IRDealID=$7

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n InvestorReportCC -c '{"Args":["SaveInvestorReport","'"$IRID"'","'"$IRData"'","'"$IRMonth"'","'"$IRYear"'","'"$IRUpdatedBy"'","'"$IRUpdationDate"'","'"$IRDealID"'"]}'
