FILES="Corefile_bloomsec_nsec
Corefile_bloomsec_nsec5
Corefile_nsec
Corefile_nsec5"

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

filename="memusage_log".$current_time".txt"
touch $filename
filename="./memusage/"$filename

cd ..

temp_za="temp_za.txt"
temp_ns="temp_ns.txt"

function conduct() {
    ./coredns -conf ./memusage/$1/$2 -dns.port 1054 | tee -a $3
}

function repeat() {
    touch $temp_za
    touch $temp_ns
    for j in {1..3}
    do 
        rm -f /memusage/testdata/db.miek.nl.signed
        conduct "Corefiles_za" $1 $temp_za
        conduct "Corefiles_ns" $2 $temp_ns
    done
    echo "================== Processing ZA $1 =================" | tee -a $filename
    grep "TotalAlloc" $temp_za | tee -a $filename
    echo "================== Processing NS $2 =================" | tee -a $filename
    grep "TotalAlloc" $temp_ns | tee -a $filename
    rm -f $temp_za
    rm -f $temp_ns
}

# main function

for f in $FILES
do
    if [[ "$f" = *"bloomsec"* ]]
    then
        for i in {1..3}
        do
            TEMP="${f}_v$i"
            repeat $TEMP $f
        done
    else 
        repeat $f $f
    fi
done



