import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_bloomsec_v1_2022.08.16-10.35.22/latencies.txt'
df=pd.read_csv(filename,sep='\\s+')
sns.displot(df, x="latency")
plt.xlim(0.0001,0.001)
plt.axvline(df.latency.mean(),
            color='red')
plt.xlabel("Latency [s]", size=12)
plt.ylabel("Number of queries", size=12)
plt.show()