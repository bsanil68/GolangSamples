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

echo "assetID is " $1
assetID=$1
echo "loanID is " $2
loanID=$2
echo "classVehicle is " $3
classVehicle=$3
echo "makersName is " $4
makersName=$4
echo "typeOfBody" $5
typeOfBody=$5
echo "horsePow is " $6
horsePow=$6
echo "chassisNum is" $7
chassisNum=$7
echo "numOfCylinders is" $8
numOfCylinders=$8
echo "yearOfManuf is" $9
yearOfManuf=$9
echo "engineNum is" ${10}
engineNum=${10}
echo "Colour is" ${11}
Colour=${11}
echo "registrationNum is" ${12}
registrationNum=${12}
echo "dateRTOGrantReg is" ${13}
dateRTOGrantReg=${13}
echo "CollateralHash is" ${14}
CollateralHash=${14}

peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n CollateralCC -c '{"Args":["CreateLoanAsset","'"$assetID"'","'"$loanID"'","'"$classVehicle"'","'"$makersName"'","'"$typeOfBody"'","'"$horsePow"'","'"$chassisNum"'","'"$numOfCylinders"'","'"$yearOfManuf"'","'"$engineNum"'","'"$Colour"'","'"$registrationNum"'","'"$dateRTOGrantReg"'","'"$CollateralHash"'"]}'
