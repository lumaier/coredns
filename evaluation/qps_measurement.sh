MAX_QPS=128000
NR_CLIENTS=24
NS_IP=129.132.85.126
PORT=1054
TIME_RUN=10

current_time=$(date "+%Y.%m.%d-%H.%M.%S")
evaldir=./measurements/qps_measurement_$1_$current_time
mkdir $evaldir

temp=$evaldir"/log.txt"
touch $temp

echo "======================= $1 ===========================" | tee -a $temp
for ((i=1000;i<=$MAX_QPS;i*=2)); do
    dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $NR_CLIENTS -Q $i -l $TIME_RUN -S 1 | tee -a $temp
done