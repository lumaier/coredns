FILES="Corefile_bloomsec_nsec
Corefile_bloomsec_nsec5
Corefile_nsec
Corefile_nsec5"

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

cd ..

evaldir=./memusage/memusage.$current_time
mkdir $evaldir
filename=$evaldir/memusage_log.txt
touch $filename

function conduct() {
    ./coredns -conf ./memusage/$1/$2 -dns.port 1054 | tee -a $3
}

function repeat() {
    temp_za=$evaldir"/temp_za."$1".txt"
    temp_ns=$evaldir"/temp_ns."$1".txt"
    touch $temp_za
    touch $temp_ns
    for j in {1..3}
    do 
        rm -f /memusage/testdata/db.miek.nl.signed
        conduct "Corefiles_za" $1 $temp_za
        conduct "Corefiles_ns" $1 $temp_ns
    done
    echo "================== Processing ZA $1 =================" | tee -a $filename
    grep "TotalAlloc" $temp_za | tee -a $filename
    echo "================== Processing NS $1 =================" | tee -a $filename
    grep "TotalAlloc" $temp_ns | tee -a $filename
}

# main function

for f in $FILES
do
    if [[ "$f" = *"bloomsec"* ]]
    then
        for i in {1..3}
        do
            TEMP="${f}_v$i"
            repeat $TEMP
        done
    else 
        repeat $f
    fi
done



