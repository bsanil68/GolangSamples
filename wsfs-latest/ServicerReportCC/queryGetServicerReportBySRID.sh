
SRID=$1

var=$(docker exec cli peer chaincode query -C mychannel -n ServicerReportCC -c '{"Args":["GetServicerReportAttributeBySRID","'"$SRID"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
