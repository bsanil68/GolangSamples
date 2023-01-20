TRID=$1
TRDealId=$2
TRNote=$3
TRCusip=$4
TROriginalBalance=$5
TRInterestRate=$6
TRUpdatedBy=$7
TRUpdationDate=$8

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n TrancheCC -c '{"Args":["SaveTranche","'"$TRID"'","'"$TRDealId"'","'"$TRNote"'","'"$TRCusip"'","'"$TROriginalBalance"'","'"$TRInterestRate"'","'"$TRUpdatedBy"'","'"$TRUpdationDate"'"]}'
