$TTL    30M
$ORIGIN miek.nl.
@       IN      SOA     ns1 miek.miek.nl. ( 1282630060 4H 1H 7D 4H )
                IN      NS      ns1
                IN      DNSKEY  257 3 13 sfzRg5nDVxbeUc51su4MzjgwpOpUwnuu81SlRHqJuXe3SOYOeypR69tZ52XLmE56TAmPHsiB8Rgk+NTpf0o1Cw==
a       IN      A       1.0.0.2
b.a     IN      A       1.0.0.3
b       IN      A       4.0.0.1
c.b     IN      A       4.0.0.2
ns1     IN      A       127.0.0.1