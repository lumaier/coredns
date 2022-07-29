NAMES="ethz.ch
google.ch
kk.migros.ch
ashdlkjafl.ch"

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

filename="messagesize".$current_time".txt"
touch $filename
temp="temp.txt"
touch $temp

host=localhost
port=53

while getopts h:p: flag
do
    case "${flag}" in
        h) host=${OPTARG};;
        p) port=${OPTARG};;
    esac
done

for name in $NAMES
do 
    dig @$host -p$port $name +dnssec | tee -a $temp
done

grep "MSG SIZE" $temp | tee -a $filename

rm -f $temp