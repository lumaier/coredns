import matplotlib.pyplot as plt
import seaborn as sns

x = ['NSEC', 'NSEC3', 'NSEC5', 'BloomSEC (fp=0.5)', 'BloomSEC (fp=0.1)', 'BloomSEC (fp=0.01)']
y = [281, 295, 435, 440, 439, 442]

sns.set(font_scale = 1.2)

g = sns.barplot(x, y)
plt.xticks(rotation=30)
plt.ylabel('seconds')
plt.show()