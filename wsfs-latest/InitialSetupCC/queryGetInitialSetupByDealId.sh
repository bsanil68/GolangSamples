ISDealId=$1

var=$(docker exec cli peer chaincode query -C mychannel -n InitialSetupCC -c '{"Args":["GetInitialSetupByDealID","'"$ISDealId"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
