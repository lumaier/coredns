import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_nsec_2022.08.16-09.09.36/latencies.txt'
df=pd.read_csv(filename,sep='\\s+')
df['System']='NSEC'

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_nsec3_2022.08.16-09.33.56/latencies.txt'
temp=pd.read_csv(filename,sep='\\s+')
temp['System']='NSEC3'

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_nsec5_2022.08.16-10.05.41/latencies.txt'
temp2=pd.read_csv(filename,sep='\\s+')
temp2['System']='NSEC5'

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_bloomsec_v1_2022.08.16-10.35.22/latencies.txt'
temp3=pd.read_csv(filename,sep='\\s+')
temp3['System']='BloomSEC (fp=0.5)'

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_bloomsec_v2_2022.08.16-11.27.52/latencies.txt'
temp4=pd.read_csv(filename,sep='\\s+')
temp4['System']='BloomSEC (fp=0.1)'

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_bloomsec_v3_2022.08.16-12.01.41/latencies.txt'
temp5=pd.read_csv(filename,sep='\\s+')
temp5['System']='BloomSEC (fp=0.001)'


df=pd.concat([df,temp,temp2,temp3,temp4,temp5],axis=0).reset_index()
# df=pd.concat([df,temp],axis=0).reset_index()
df['latency'] = df['latency']*1000.0

sns.set(font_scale = 1.2)
g = sns.displot(x='latency', hue='System', data=df)
plt.xlabel("Latency [ms]")
plt.ylabel("Number of queries")

plt.xlim(0.15,0.5)
plt.grid()
plt.show()