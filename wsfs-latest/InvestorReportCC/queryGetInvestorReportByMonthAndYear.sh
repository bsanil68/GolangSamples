IRDealID=$1
IRMonth=$2
IRYear=$3

var=$(docker exec cli peer chaincode query -C mychannel -n InvestorReportCC -c '{"Args":["GetInvestorReportByMonthAndYear","'$IRDealID'","'"$IRMonth"'","'"$IRYear"'"]}')
echo $var | sed -e 's/Query\ Result\:\ //g'
