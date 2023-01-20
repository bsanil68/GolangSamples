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

echo "PoolId is:" $1
PoolId=$1
echo "PoolName is:" $2
PoolName=$2
echo "PoolDesc is:" $3
PoolDesc=$3
echo "PoolOwner is:" $4
PoolOwner=$4
echo "PoolStartedDate is:" $5
PoolStartedDate=$5
echo "PoolExpiryDate is:" $6
PoolExpiryDate=$6
echo "PoolSts is:" $7
PoolSts=$7
echo "ApprovalDate is:" $8
ApprovalDate=$8
echo "ApprovalId is:" $9
ApprovalId=$9
echo "NumAssets is:" ${10}
NumAssets=${10}
echo "PoolHash is:" ${11}
PoolHash=${11}
echo "PoolCreatedDate is:" ${12}
PoolCreatedDate=${12}
echo "WarehouseLenderName is:" ${13}
WarehouseLenderName=${13}
peer chaincode invoke -o orderer0:7050  --tls --cafile /root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg0/orderers/orderer0.ordererorg0/msp/tlscacerts/tlsca.ordererorg0-cert.pem -C composerchannelrest -n WarehouseCreatepoolCC -v 1 -c '{"Args":["Createpool","'"$PoolId"'","'"$PoolName"'","'"$PoolDesc"'","'"$PoolOwner"'","'"$PoolStartedDate"'","'"$PoolExpiryDate"'","'"$PoolSts"'","'"$ApprovalDate"'","'"$ApprovalId"'","'"$NumAssets"'","'"$PoolHash"'","'"$PoolCreatedDate"'","'"$WarehouseLenderName"'"]}'


