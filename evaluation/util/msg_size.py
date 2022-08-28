import matplotlib.pyplot as plt
import seaborn as sns

x = ['NSEC', 'NSEC3', 'NSEC5', 'BloomSEC (fp=0.5)', 'BloomSEC (fp=0.1)', 'BloomSEC (fp=0.01)']
y = [564, 634, 1194, 1322, 1424.4, 1447.44]



sns.set(font_scale = 1.2)

g = sns.barplot(x, y)
plt.axhline(y=1450,color='r',linestyle='--',label='MTU')
plt.xticks(rotation=30)
plt.ylabel('bytes')
plt.text(x=0.05,y=1350, s="MTU", c='r')
plt.show()