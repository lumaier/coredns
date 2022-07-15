$TTL    30M
$ORIGIN miek.nl.
@   IN  SOA	linode.atoom.net. miek.miek.nl. (1282630060 4H 1H 7D 4H)
			1800	NS	ext.ns.whyscream.net.
			1800	NS	omval.tednet.nl.
			1800	NS	linode.atoom.net.
			1800	NS	ns-ext.nlnetlabs.nl.
			1800	A	139.162.196.78
			1800	MX	1 aspmx.l.google.com.
			1800	MX	5 alt1.aspmx.l.google.com.
			1800	MX	5 alt2.aspmx.l.google.com.
			1800	MX	10 aspmx2.googlemail.com.
			1800	MX	10 aspmx3.googlemail.com.
			1800	AAAA	2a01:7e00::f03c:91ff:fef1:6735
			1800	DNSKEY	256 3 8 (
					AwEAAcNEU67LJI5GEgF9QLNqLO1SMq1EdoQ6
					E9f85ha0k0ewQGCblyW2836GiVsm6k8Kr5EC
					IoMJ6fZWf3CQSQ9ycWfTyOHfmI3eQ/1Covhb
					2y4bAmL/07PhrL7ozWBW3wBfM335Ft9xjtXH
					Py7ztCbV9qZ4TVDTW/Iyg0PiwgoXVesz
					) ; ZSK; alg = RSASHA256; key id = 12051
			1800	DNSKEY	257 3 8 (
					AwEAAcWdjBl4W4wh/hPxMDcBytmNCvEngIgB
					9Ut3C2+QI0oVz78/WK9KPoQF7B74JQ/mjO4f
					vIncBmPp6mFNxs9/WQX0IXf7oKviEVOXLjct
					R4D1KQLX0wprvtUIsQFIGdXaO6suTT5eDbSd
					6tTwu5xIkGkDmQhhH8OQydoEuCwV245ZwF/8
					AIsqBYDNQtQ6zhd6jDC+uZJXg/9LuPOxFHbi
					MTjp6j3CCW0kHbfM/YHZErWWtjPj3U3Z7knQ
					SIm5PO5FRKBEYDdr5UxWJ/1/20SrzI3iztvP
					wHDsA2rdHm/4YRzq7CvG4N0t9ac/T0a0Sxba
					/BUX2UVPWaIVBdTRBtgHi0s=
					) ; KSK; alg = RSASHA256; key id = 33694

a.miek.nl.		1800	IN A	139.162.196.78
				1800	AAAA	2a01:7e00::f03c:91ff:fef1:6735

archive.miek.nl.	1800	IN CNAME a.miek.nl.

go.dns.miek.nl.		1800	IN TXT	"Hello!"

www.miek.nl.		1800	IN CNAME a.miek.nl.