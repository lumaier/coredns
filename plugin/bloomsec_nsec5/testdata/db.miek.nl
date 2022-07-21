$TTL    30M
$ORIGIN miek.nl.
@       IN      SOA     ns1 miek.miek.nl. ( 1282630060 4H 1H 7D 4H )
                IN      NS      ns1
                IN      A       1.0.0.1

a               IN      A       1.0.0.2
b.a             IN      A       1.0.0.3
b               IN      A       4.0.0.1
c.b             IN      A       4.0.0.2
ns1             IN      A       127.0.0.1