dir=./measurements/qps_measurement_nsec_2022.08.16-08.19.18

file=$dir/qps.txt
out=$dir/pretty_qps.txt

grep 'Queries per second' $file | cut -c25- | tee -a $out