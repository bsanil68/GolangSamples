IRID=$1

var=$(docker exec cli peer chaincode query -C mychannel -n InvestorReportCC -c '{"Args":["GetInvestorReportByID","'"$IRID"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
