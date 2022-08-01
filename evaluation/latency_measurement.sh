NR_CLIENTS=24
NS_IP=127.0.0.1
PORT=1054
TIME_RUN=1

current_time=$(date "+%Y.%m.%d-%H.%M.%S")
evaldir=./measurements/latency_measurement_$1_$current_time
mkdir $evaldir

log=$evaldir"/log.txt"
touch $log

echo "======================= $1 ===========================" | tee -a $log
dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $NR_CLIENTS -Q 1000 -l $TIME_RUN -v | tee -a $log

filtered=$evaldir"/latencies.txt"
touch $filtered

echo "latency" | tee -a $filtered
grep "NXDOMAIN" $log | cut -c33- | tee -a $filtered
sed -i '$ d' $filtered





