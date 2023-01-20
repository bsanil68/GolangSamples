export CORE_PEER_LOCALMSPID="ServicerMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/peers/peer1.servicer/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/users/Admin@servicer/msp
export CORE_PEER_ADDRESS=peer1:7051
export  CORE_PEER_TLS_CLIENTROOTCAS_FILES=/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/users/Admin@servicer/tls/ca.crt
export  CORE_PEER_TLS_CLIENTCERT_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/users/Admin@servicer/tls/client.crt
export CORE_PEER_TLS_CLIENTKEY_FILE=/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/users/Admin@servicer/tls/client.key
export CORE_PEER_TLS_ENABLED=true

export ORDERER_GENERAL_TLS_PRIVATEKEY=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg/orderers/orderer.ordererorg/tls/server.key
export ORDERER_GENERAL_TLS_CERTIFICATE=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg/orderers/orderer.ordererorg/tls/server.crt
export ORDERER_GENERAL_TLS_ROOTCAS=[/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg/orderers/orderer.ordererorg/tls/ca.crt,/root/bcnetwork/conf/crypto-config/peerOrganizations/servicer/peers/peer1.servicer/tls/ca.crt]

export $CORE_PEER_TLS_ENABLED=/root/bcnetwork/conf/crypto-config/ordererOrganizations/ordererorg/orderers/orderer.ordererorg/tls/ca.crt

peer chaincode install -n DLoanMapCC -v 1 -p github.com/intainabs/chaincode/DLoanMapCC
