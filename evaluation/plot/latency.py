import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
filename='/home/luca/Documents/bth/coredns/evaluation/measurements/latency_measurement_bloomsec_nsec5_2022.08.01-14.34.54/latencies.txt'
df=pd.read_csv(filename,sep='\\s+')
sns.displot(df, x="latency")
plt.show()