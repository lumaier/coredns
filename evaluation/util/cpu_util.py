import matplotlib.pyplot as plt
import seaborn as sns

x = ['NSEC', 'NSEC3', 'NSEC5', 'BloomSEC (fp=0.5)', 'BloomSEC (fp=0.1)', 'BloomSEC (fp=0.01)']
y = [50, 55, 120, 100, 60, 52]

sns.set(font_scale = 1.2)

g = sns.barplot(x, y)
plt.xticks(rotation=30)
plt.ylabel('% CPU utilization')
plt.show()