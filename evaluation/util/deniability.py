from math import log2
import numpy as np
import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd

fp1 = 0.01
fp2 = 0.1
fp3 = 0.5

n = 150000000.

x = np.arange(150000000,10**11,100000)


y1= (1-4**(-(x-n)*fp1/n ))**(-log2(fp1))
y2= (1-4**(-(x-n)*fp2/n ))**(-log2(fp2))
y3= (1-4**(-(x-n)*fp3/n ))**(-log2(fp3))

x_df = pd.DataFrame({"numbers": x})
df1 = pd.DataFrame({"BloomSEC (fp=0.01)": y1})
df2 = pd.DataFrame({"BloomSEC (fp=0.1)": y2})
df3 = pd.DataFrame({"BloomSEC (fp=0.5)": y3})

df = pd.concat([x_df,df1,df2,df3],axis=1)
dfm = pd.melt(df, id_vars='numbers', var_name='System', value_name='vals')

r=sns.lineplot(data=dfm, x="numbers", y="vals", hue="System", linewidth = 1,legend='brief')
r.set_ylabel("$\gamma$-deniability")
r.set_xlabel("size of the dictionary")

plt.ylim(0,1)
plt.xlim(150000000,10**11)
plt.xscale('log')

# Show the plot
plt.legend()
plt.grid()
plt.show()