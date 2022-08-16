MAX_QPS=96000
NR_CLIENTS=24
NS_IP=167.71.45.171
# NS_IP=127.0.0.1
PORT=1054
TIME_RUN=10

current_time=$(date "+%Y.%m.%d-%H.%M.%S")
evaldir=./measurements/qps_measurement_$1_$current_time
mkdir $evaldir

temp=$evaldir"/log.txt"
touch $temp

echo "======================= $1 ===========================" | tee -a $temp
for ((i=2000;i<=$MAX_QPS;i+=2000)); do
    dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $NR_CLIENTS -Q $i -l $TIME_RUN -S 1 | tee -a $temp
done

filtered=$evaldir"/qps.txt"
touch $filtered

out=$evaldir"/pretty_qps.txt"
touch $out

grep 'Queries per second\|Average Latency\|Latency StdDev\|Response codes' $temp | tee -a $filtered

grep 'Queries per second' $filtered | cut -c25- | tee -a $out