MAX_QPS=64000
NR_CLIENTS=24
NS_IP=167.71.45.171
# NS_IP=127.0.0.1
PORT=1054
TIME_RUN=10

qps=(1000
2000
4000
8000
12000
16000
20000
24000
28000
32000
40000
64000
96000
128000)

current_time=$(date "+%Y.%m.%d-%H.%M.%S")
evaldir=./measurements/qps_measurement_$1_$current_time
mkdir $evaldir

temp=$evaldir"/log.txt"
touch $temp

echo "======================= $1 ===========================" | tee -a $temp
for i in "${qps[@]}"; do
    dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $NR_CLIENTS -Q $i -l $TIME_RUN -S 1 | tee -a $temp
done

filtered=$evaldir"/qps.txt"
touch $filtered

grep 'Queries per second\|Average Latency\|Latency StdDev\|Response codes' $temp | tee -a $filtered