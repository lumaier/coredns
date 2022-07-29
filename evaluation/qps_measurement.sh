MAX_QPS=128000
MAX_CLIENTS=24
NS_IP=127.0.0.1
PORT=1054

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

temp="temp_$1_$current_time.txt"
touch $temp

for ((c=4;c<=$MAX_CLIENTS;c+=4)); do
    echo "=========== For $c clients: ===============" | tee -a $temp
    for ((i=1000;i<=$MAX_QPS;i*=2)); do
        dnsperf -d ./testdata/queries/queries_nex -D -p $PORT -s $NS_IP -c $c -Q $i -l 10 -S 1 | tee -a $temp
    done
done