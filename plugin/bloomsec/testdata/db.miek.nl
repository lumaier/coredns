$TTL    30M
$ORIGIN miek.nl.
@       IN      SOA     ns1 miek.miek.nl. ( 1282630060 4H 1H 7D 4H )
                IN      NS      ns1
                IN      MX      1  aspmx.l.google.com.
		IN 	AAAA    2a01:7e00::f03c:91ff:fe79:234c

a		IN 	A       1.0.0.2
www     	IN 	CNAME 	a
b.a             IN      A       1.0.0.3
b               IN      A       4.0.0.1
c.b             IN      A       4.0.0.2

bla                 IN  NS      ns1.bla.com.
ns3.blaaat.miek.nl. IN  AAAA    ::1 ; non-glue, should be signed.
; in baliwick nameserver that requires glue, should not be signed
bla                 IN  NS      n1.miek.nl
ns1                 IN  A       127.0.0.1