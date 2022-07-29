temp1="temp1.txt"
temp2="temp2.txt"

touch $temp1
touch $temp2

for i in 1000; do
    dnsperf -d ./testdata/queries/queries_nex -D -p 1054 -s 127.0.0.1 -c 24 -Q $i -l 10 -v | tee -a $temp1
    dnsperf -d ./testdata/queries/queries_nex -D -p 1054 -s 127.0.0.1 -c 24 -Q $i -l 10 -S 1 | tee -a $temp2
done