TRID=$1

var=$(docker exec cli peer chaincode query -C mychannel -n TrancheCC -c '{"Args":["GetTrancheByID","'"$TRID"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
