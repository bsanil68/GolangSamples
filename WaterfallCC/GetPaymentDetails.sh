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

poolId=$1

month=$2

if [ -e getPaymentDetails.json ]
then
    rm -rf getPaymentDetails.json
fi
var=$(peer chaincode query -C composerchannelrest -n Waterfall5CC -c '{"Args":["GetPaymentWaterfallDetails","'"$poolId"'","'"$month"'"]}')
echo $var | sed -e 's/Query\ Result\://g'>>getPaymentDetails.json
