NR_CLIENTS=24
NS_IP=167.71.45.171
PORT=1054
TIME_RUN=120

dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $NR_CLIENTS -Q 4000 -l $TIME_RUN -v 





