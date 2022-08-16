import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np

myRange = np.arange(2000,96001,2000) # Could be something like myRange = range(1,1000,1)
df = pd.DataFrame({"numbers": myRange})


filename='/home/luca/Documents/bth/coredns/evaluation/measurements/qps_measurement_bloomsec_v1_2022.08.16-10.26.31/pretty_qps.txt'
df1=pd.read_csv(filename,sep='\\s+',names=['BloomSEC (fp=0.5)'])

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/qps_measurement_nsec3_2022.08.16-09.24.51/pretty_qps.txt'
df2=pd.read_csv(filename,sep='\\s+',names=['NSEC3'])

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/qps_measurement_nsec5_2022.08.16-09.56.30/pretty_qps.txt'
df4=pd.read_csv(filename,sep='\\s+',names=['NSEC5'])

filename='/home/luca/Documents/bth/coredns/evaluation/measurements/qps_measurement_nsec_2022.08.16-08.53.38/pretty_qps.txt'
df3=pd.read_csv(filename,sep='\\s+',names=['NSEC'])

df = pd.concat([df,df1,df2,df3,df4],axis=1)
dfm = pd.melt(df, id_vars='numbers', var_name='System', value_name='vals')

r=sns.lineplot(data=dfm, x="numbers", y="vals", hue="System", linewidth = 1,legend='brief')
r.set(yscale='log')
r.set_ylabel("achieved throughput", fontsize = 12)
r.set_xlabel("queries/second", fontsize = 12)

plt.grid()
plt.show()
