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

echo "poolId is " $1
poolId=$1

echo "trancheId is " $2
trancheId=$2

echo "className is " $3
className=$3

echo "percentOfCollateralValue is " $4
percentOfCollateralValue=$4

echo "nominalAmount is " $5
nominalAmount=$5

echo "issuePricePercent is " $6
issuePricePercent=$6

echo "size is " $7
size=$7

echo "interestRateType is " $8
interestRateType=$8

echo "couponRate is " $9
couponRate=$9

echo "interestFrequency is " ${10}
interestFrequency=${10}

echo "faceValue is " ${11}
faceValue=${11}

echo "minUnitsOfSubscription is " ${12}
minUnitsOfSubscription=${12}

echo "totalUnitsAvailable is " ${13}
totalUnitsAvailable=${13}

echo "notePeriod is " ${14}
notePeriod=${14}

echo "noOfUnitsLeft is " ${15}
noOfUnitsLeft=${15}

echo "Currency " ${16}
Currency=${16}


peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n TrancheCC -c '{"Args":["SaveTranche","'"$poolId"'","'"$trancheId"'","'"$className"'","'"$percentOfCollateralValue"'","'"$nominalAmount"'","'"$issuePricePercent"'","'"$size"'","'"$interestRateType"'","'"$couponRate"'","'"$interestFrequency"'","'"$faceValue"'","'"$minUnitsOfSubscription"'","'"$totalUnitsAvailable"'","'"$notePeriod"'","'"$noOfUnitsLeft"'","'"$Currency"'"]}'

