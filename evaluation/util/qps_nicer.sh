cd ../measurements

prefix="[Status] Command line: dnsperf -d ./testdata/queries/queries_nex -D -p 1054 -s 127.0.0.1 -c 24 -Q" 

combined=combined.txt
rm -f $combined
touch $combined

for f in *; do
    if [ -d "$f" ] && [[ $f == *"qps"* ]]; then
        cd $f
        file=data.txt
        rm $file
        touch $file
        grep "Queries per second" log.txt | cut -c25- | tee -a $file
        cd ..
    fi
done