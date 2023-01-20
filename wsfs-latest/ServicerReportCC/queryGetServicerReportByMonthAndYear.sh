SRDealID=$1
SRMonth=$2
SRYear=$3

var=$(docker exec cli peer chaincode query -C mychannel -n ServicerReportCC -c '{"Args":["GetServicerReportAttributeByMonthAndYear","'"$SRDealID"'","'"$SRMonth"'","'"$SRYear"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
