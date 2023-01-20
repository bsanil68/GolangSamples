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

echo "CollectionID is:" $1
CollectionID=$1

echo "LoanID is:" $2
LoanID=$2

echo "InstallmentNumber is:" $3
InstallmentNumber=$3

echo "EmiPaid is:" $4
EmiPaid=$4

echo "InterestAmountRepaid is:" $5
InterestAmountRepaid=$5

echo "PrincipalAmountRepaid is:" $6
PrincipalAmountRepaid=$6

echo "OutstandingPrincipalBalance is:" $7
OutstandingPrincipalBalance=$7

echo "OverdueEMINumbers1To30Days is:" $8
OverdueEMINumbers1To30Days=$8

echo "OverdueEMINumbers31To60Days is:" $9
OverdueEMINumbers31To60Days=$9

echo "OverdueEMINumbers61To90Days is:" ${10}
OverdueEMINumbers61To90Days=${10}

echo "OverdueEMINumbers91To120Days is:" ${11}
OverdueEMINumbers91To120Days=${11}

echo "OverdueEMINumbers121To180Days is:" ${12}
OverdueEMINumbers121To180Days=${12}

echo "OverdueEMINumbers180PlusDays is:" ${13}
OverdueEMINumbers180PlusDays=${13}

echo "Month is:" ${14}
Month=${14}

echo "Currency is:" ${15}
Currency=${15}

echo "Collection Hash is:" ${16}
CollectionHash=${16}

peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n Collection3CC -c '{"Args":["SaveCollection","'"$CollectionID"'","'"$LoanID"'","'"$InstallmentNumber"'","'"$EmiPaid"'","'"$InterestAmountRepaid"'","'"$PrincipalAmountRepaid"'","'"$OutstandingPrincipalBalance"'","'"$OverdueEMINumbers1To30Days"'","'"$OverdueEMINumbers31To60Days"'","'"$OverdueEMINumbers61To90Days"'","'"$OverdueEMINumbers91To120Days"'",
"'"$OverdueEMINumbers121To180Days"'", "'"$OverdueEMINumbers180PlusDays"'", "'"$Month"'", "'"$Currency"'", "'"$CollectionHash"'"]}'
