NAMES="a
b
c
d
e
f
g
h
i
j
k
l
m
n
o
p"

current_time=$(date "+%Y.%m.%d-%H.%M.%S")

dir="./measurements/messagesize_"$1_$current_time
mkdir $dir
filename=$dir"/msgsize.txt"
touch $filename
temp=$dir"/log.txt"
touch $temp

host=localhost
port=1054

while getopts h:p: flag
do
    case "${flag}" in
        h) host=${OPTARG};;
        p) port=${OPTARG};;
    esac
done

for name in $NAMES
do 
    dig @$host -p$port $name".miek.nl" +dnssec | tee -a $temp
done

grep "MSG SIZE" $temp | tee -a $filename