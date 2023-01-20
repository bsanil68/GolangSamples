SRID=$1
SRKey=$2
SRValue=$3
SRMonth=$4
SRYear=$5
SRUpdatedBy=$6
SRUpdationDate=$7
SRSeqNum=$8
SRDealID=$9

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n ServicerReportCC -c '{"Args":["SaveServicerReportAttribute","'"$SRID"'","'"$SRKey"'","'"$SRValue"'","'"$SRMonth"'","'"$SRYear"'","'"$SRUpdatedBy"'","'"$SRUpdationDate"'","'"$SRSeqNum"'","'"$SRDealID"'"]}'
