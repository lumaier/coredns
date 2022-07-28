# need to specify
homedir="/home/luca/Documents/bth/coredns"

FILES="Corefile_bloomsec_nsec
Corefile_bloomsec_nsec5
Corefile_nsec
Corefile_nsec5"

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

filename="memusage_log".$current_time".txt"
touch $filename

cd $homedir

filename="./memusage/"$filename

function conduct() {
    echo "" | tee -a $filename
    echo "====================== Processing $2 =================================" | tee -a $filename
    echo "" | tee -a $filename
    echo "" | tee -a $filename
    ./coredns -conf ./memusage/$1/$2 -dns.port 1054 | tee -a $filename
    echo ""
}

# main function

for f in $FILES
do
    if [[ "$f" = *"bloomsec"* ]]
    then
        for i in {1..3}
        do
            rm -f /memusage/testdata/db.miek.nl.signed
            TEMP="${f}_v$i"
            conduct "Corefiles_za" $TEMP
            conduct "Corefiles_ns" $f
        done
    else 
        rm -f /memusage/testdata/db.miek.nl.signed
        conduct "Corefiles_za" $f
        conduct "Corefiles_ns" $f
    fi
done