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

echo "trancheID is" $1
trancheID=$1

echo "investorID is " $2
investorID=$2

echo "numOfUnitsPurchased is" $3
numOfUnitsPurchased=$3

echo "amountPaidByInvestor is" $4
amountPaidByInvestor=$4

echo "applicationDate is" $5
applicationDate=$5

echo "allocationDate is" $6
allocationDate=$6

echo "monthlyInterestForInvestor is" $7
monthlyInterestForInvestor=$7

echo "monthlyPrincipalForInvestor is" $8
monthlyPrincipalForInvestor=$8

echo "numOfInstallments is" $9
numOfInstallments=$9

echo "Currency is" ${10}
Currency=${10}
 
peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n TrancheInvestorMap8CC -c '{"Args":["SaveTrancheInvestorMap","'"$trancheID"'","'"$investorID"'","'"$numOfUnitsPurchased"'","'"$amountPaidByInvestor"'","'"$applicationDate"'","'"$allocationDate"'","'"$monthlyInterestForInvestor"'","'"$monthlyPrincipalForInvestor"'","'"$numOfInstallments"'","'"$Currency"'"]}'

