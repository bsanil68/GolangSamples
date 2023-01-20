TRDealId=$1

var=$(docker exec cli peer chaincode query -C mychannel -n TrancheCC -c '{"Args":["GetTrancheByDealID","'"$TRDealId"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
