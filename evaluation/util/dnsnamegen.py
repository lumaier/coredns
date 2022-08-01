import sys
import string
import random

def main():
    args = sys.argv[1:]
    n = int(args[0])                    # number of records
    positive_queries = int(args[1])     # number of positive queries
    total_queries = int(args[2])        # number of total queries
    # flag = int(args[3])                 # 0: random domain names with 1 label of length 10
    #                                     # 1: random domain names with 3 labels of length 2
    

    f = open("db.miek.nl", "w")
    f_queries = open("queries", "w")
    f_nex_queries = open("queries_nex", "w")

    f.write("$TTL    30M\n$ORIGIN miek.nl.\n@       IN      SOA     ns1 miek.miek.nl. ( 1282630060 4H 1H 7D 4H )\n")
    
    for i in range(n):
        a = randStr()
        f.write(a + '\tIN\tA\t' + randIP()+'\n')
        if i < positive_queries:
            f_queries.write(a + ".miek.nl" + "\tA\n")

    for i in range(total_queries - positive_queries):
        f_queries.write(randStr() + ".miek.nl\tA\n")

    for i in range(total_queries):
        f_nex_queries.write(randStr() + ".miek.nl\tA\n")

    f.close()
    f_queries.close()
    f_nex_queries.close()

def randStr(chars = string.ascii_lowercase + string.digits, N=10):
	return "".join(random.choice(chars) for _ in range(N))
    
def randIP(chars = string.digits, N=10):
	return '10.' + str(random.randint(0,253)) + '.' + str(random.randint(0,253)) + '.' + str(random.randint(0,253))

if __name__== "__main__":
  main()