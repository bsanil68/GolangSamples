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

echo "emiID is:" $2
emiID=$2

echo "InstallmentNumber is:" $3
InstallmentNumber=$3

echo "ScheduleDate:" $4
ScheduleDate=$4

echo "EmiAmount:" $5
EmiAmount=$5

echo "installPrinAmount:" $6
installPrinAmount=$6

echo "installIntAmt:" $7
installIntAmt=$7

echo "outstandPrincp:" $8
outstandPrincp=$8

echo "Currency:" $9
Currency=$9

echo "EMIHash:" ${10}
EMIHash=${10}

peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n EMIScheduleCC -c '{"Args":["CreateEMISchedule","'"$LoanID"'","'"$emiID"'","'"$InstallmentNumber"'","'"$ScheduleDate"'","'"$EmiAmount"'","'"$installPrinAmount"'","'"$installIntAmt"'","'"$outstandPrincp"'","'"$Currency"'","'"$EMIHash"'"]}'



