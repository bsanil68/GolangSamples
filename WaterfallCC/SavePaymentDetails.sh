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

echo "transactionID is" $1
transactionID=$1

echo "poolId is" $2
poolId=$2

echo "taxPayment is" $3
taxPayment=$3

echo "paymentOfReimbursableExpenses is" $4
paymentOfReimbursableExpenses=$4

echo "paymentOfServicingFeesToServicer is" $5
paymentOfServicingFeesToServicer=$5

echo "interestPayableToClassA is " $6
interestPayableToClassA=$6

echo "principalPayableToClassA is " $7
principalPayableToClassA=$7

echo "interestPayableToClassB is " $8
interestPayableToClassB=$8

echo "principalPayableToClassB is " $9
principalPayableToClassB=$9

echo "interestPayableToClassC is " ${10}
interestPayableToClassC=${10}

echo "principalPayableToClassC is " ${11}
principalPayableToClassC=${11}

echo "balanceAsIncomeOfClassC is " ${12}
balanceAsIncomeOfClassC=${12}

echo "month is " ${13}
month=${13}

echo "totalCollections is " ${14}
totalCollections=${14}

echo "Currency is " ${15}
Currency=${15}

peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n Waterfall5CC -c '{"Args":["SavePaymentWaterfallDetails","'"$transactionID"'","'"$poolId"'","'"$taxPayment"'","'"$paymentOfReimbursableExpenses"'","'"$paymentOfServicingFeesToServicer"'","'"$interestPayableToClassA"'","'"$principalPayableToClassA"'","'"$interestPayableToClassB"'","'"$principalPayableToClassB"'","'"$interestPayableToClassC"'","'"$principalPayableToClassC"'","'"$balanceAsIncomeOfClassC"'","'"$month"'","'"$totalCollections"'","'"$Currency"' "]}'

