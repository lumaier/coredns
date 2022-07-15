package bloomfile_nsec5

import (
	"context"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

// All OPT RR are added in server.go, so we don't specify them in the unit tests.
var dnssecTestCases = []test.Case{
	{
		Qname: "miek.nl.", Qtype: dns.TypeSOA, Do: true,
		Answer: []dns.RR{
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
		Ns: auth,
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeAAAA, Do: true,
		Answer: []dns.RR{
			test.AAAA("miek.nl.	1800	IN	AAAA	2a01:7e00::f03c:91ff:fef1:6735"),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	AAAA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. oE5BbhEnBdGxWXKqt102RAce/Ef/PllZXXu2jPImKw1ouAAdcAXa4uGJAZVfhPSHM4PvtoGIxifgb0bXBc6TJw=="),
		},
		Ns: auth,
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeNS, Do: true,
		Answer: []dns.RR{
			test.NS("miek.nl.	1800	IN	NS	ext.ns.whyscream.net."),
			test.NS("miek.nl.	1800	IN	NS	linode.atoom.net."),
			test.NS("miek.nl.	1800	IN	NS	ns-ext.nlnetlabs.nl."),
			test.NS("miek.nl.	1800	IN	NS	omval.tednet.nl."),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	NS 13 2 1800 20220723181226 20220621032053 59725 miek.nl. Wy/bvt7X8qZ1ZOQsHJk+K568VNiv30XvG01h5jq2lvOwh4xc+fivLlOhi8rAV0jduuAUZdu50SSMlQ3XrXYE/Q=="),
		},
	},
	{
		Qname: "miek.nl.", Qtype: dns.TypeMX, Do: true,
		Answer: []dns.RR{
			test.MX("miek.nl.	1800	IN	MX	1 aspmx.l.google.com."),
			test.MX("miek.nl.	1800	IN	MX	10 aspmx2.googlemail.com."),
			test.MX("miek.nl.	1800	IN	MX	10 aspmx3.googlemail.com."),
			test.MX("miek.nl.	1800	IN	MX	5 alt1.aspmx.l.google.com."),
			test.MX("miek.nl.	1800	IN	MX	5 alt2.aspmx.l.google.com."),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	MX 13 2 1800 20220723181226 20220621032053 59725 miek.nl. m21ROZl3a4oiQCP2nf/HmWvrCv99ynVw1ngJvMSOCj80X1zAV64c4v5ZC9J4ARF6l3T2KlCzsgior8vG1vSJdA=="),
		},
		Ns: auth,
	},
	{
		Qname: "www.miek.nl.", Qtype: dns.TypeA, Do: true,
		Answer: []dns.RR{
			test.A("a.miek.nl.	1800	IN	A	139.162.196.78"),
			test.RRSIG("a.miek.nl.	1800	IN	RRSIG	A 8 3 1800 20160426031301 20160327031301 12051 miek.nl. lxLotCjWZ3kihTxk="),
			test.CNAME("www.miek.nl.	1800	IN	CNAME	a.miek.nl."),
			test.RRSIG("www.miek.nl.	1800	IN	RRSIG	CNAME 13 3 1800 20220723192909 20220621054528 59725 miek.nl. Kw3EO7EbeZWTCzLRMTkT2A0auETX/prSavCg6AuAr3C9iTitoy0kI7AG3SJ/FBYJk64Fke7ax6oPE8f1po1G8g=="),
		},
		Ns: auth,
	},
	{
		// NoData
		Qname: "a.miek.nl.", Qtype: dns.TypeSRV, Do: true,
		Ns: []dns.RR{
			test.NSEC("a.miek.nl.	14400	IN	NSEC	archive.miek.nl. A AAAA RRSIG NSEC"),
			test.RRSIG("a.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. aNBNT2VX17Pa716Gi+Fc/H4ZeSrjsCR53U44EJz9EIBYMA1urlhEeGi7g0H9GiyVxKOjbBSc51PWatV+UCUBRA=="),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
	},
	{
		Qname: "b.miek.nl.", Qtype: dns.TypeA, Do: true,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.RRSIG("_bf1.miek.nl.	14400	IN	RRSIG	TXT 13 3 14400 20220723181226 20220621032053 59725 miek.nl. QUe5qipncsgjDTrgQ/9DLQCyw0YM6qLQWnesVchJueQMSSbXELDikaqb7naD80SF51Jjpkr8REnUZqfe0o0MXA=="),
			test.TXT("_bf1.miek.nl.	14400	IN	TXT	\"arZMWIcKXUNN\" \"144\" \"5\""),
			test.NSEC("archive.miek.nl.	14400	IN	NSEC	go.dns.miek.nl. CNAME RRSIG NSEC"),
			test.RRSIG("archive.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. Sqdlq6upOFmgFUmCQ1hUNRCX41qwHhENHjTEk2ctLgRoJIGVI6KP3P/viqQYdcesHSxO8xo8NV3F9NZoPSi2JA=="),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
	},
	{
		Qname: "b.blaat.miek.nl.", Qtype: dns.TypeA, Do: true,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.RRSIG("_bf1.miek.nl.	14400	IN	RRSIG	TXT 13 3 14400 20220723181226 20220621032053 59725 miek.nl. QUe5qipncsgjDTrgQ/9DLQCyw0YM6qLQWnesVchJueQMSSbXELDikaqb7naD80SF51Jjpkr8REnUZqfe0o0MXA=="),
			test.TXT("_bf1.miek.nl.	14400	IN	TXT	\"arZMWIcKXUNN\" \"144\" \"5\""),
			test.NSEC("archive.miek.nl.	14400	IN	NSEC	go.dns.miek.nl. CNAME RRSIG NSEC"),
			test.RRSIG("archive.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. Sqdlq6upOFmgFUmCQ1hUNRCX41qwHhENHjTEk2ctLgRoJIGVI6KP3P/viqQYdcesHSxO8xo8NV3F9NZoPSi2JA=="),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
	},
	{
		Qname: "b.a.miek.nl.", Qtype: dns.TypeA, Do: true,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			// dedupped NSEC, because 1 nsec tells all
			test.NSEC("a.miek.nl.	14400	IN	NSEC	archive.miek.nl. A AAAA RRSIG NSEC"),
			test.RRSIG("a.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. aNBNT2VX17Pa716Gi+Fc/H4ZeSrjsCR53U44EJz9EIBYMA1urlhEeGi7g0H9GiyVxKOjbBSc51PWatV+UCUBRA=="),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
	},
}

var auth = []dns.RR{
	test.NS("miek.nl.	1800	IN	NS	ext.ns.whyscream.net."),
	test.NS("miek.nl.	1800	IN	NS	linode.atoom.net."),
	test.NS("miek.nl.	1800	IN	NS	ns-ext.nlnetlabs.nl."),
	test.NS("miek.nl.	1800	IN	NS	omval.tednet.nl."),
	test.RRSIG("miek.nl.	1800	IN	RRSIG	NS 13 2 1800 20220723181226 20220621032053 59725 miek.nl. Wy/bvt7X8qZ1ZOQsHJk+K568VNiv30XvG01h5jq2lvOwh4xc+fivLlOhi8rAV0jduuAUZdu50SSMlQ3XrXYE/Q=="),
}

func TestLookupDNSSEC(t *testing.T) {
	zone, err := Parse(strings.NewReader(dbMiekNLSigned), testzone, "stdin", 0)
	if err != nil {
		t.Fatalf("Expected no error when reading zone, got %q", err)
	}

	fm := File{Next: test.ErrorHandler(), Zones: Zones{Z: map[string]*Zone{testzone: zone}, Names: []string{testzone}}}
	ctx := context.TODO()

	for _, tc := range dnssecTestCases {
		m := tc.Msg()

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := fm.ServeDNS(ctx, rec, m)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		resp := rec.Msg
		if err := test.SortAndCheck(resp, tc); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkFileLookupDNSSEC(b *testing.B) {
	zone, err := Parse(strings.NewReader(dbMiekNLSigned), testzone, "stdin", 0)
	if err != nil {
		return
	}

	fm := File{Next: test.ErrorHandler(), Zones: Zones{Z: map[string]*Zone{testzone: zone}, Names: []string{testzone}}}
	ctx := context.TODO()
	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	tc := test.Case{
		Qname: "b.miek.nl.", Qtype: dns.TypeA, Do: true,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.TXT("_bf0.miek.nl.	14400	IN	TXT	\"iE0SZVIQp2ZV\" \"144\" \"5\""),
			test.RRSIG("_bf0.miek.nl.	14400	IN	RRSIG	TXT 13 3 14400 20220723181226 20220621032053 59725 miek.nl. zChRSOfjSg+adtj6M1Of7B4seFHVJB+OQ1O7m5ob9cvZhLi+2ADAP6cd7JKDrzOd1LQnyMEK1Vvy+OpI222MfQ=="),
			test.NSEC("miek.nl.	14400	IN	NSEC	a.miek.nl. A NS SOA MX AAAA RRSIG NSEC DNSKEY"),
			test.RRSIG("miek.nl.	14400	IN	RRSIG	NSEC 13 2 14400 20220723181226 20220621032053 59725 miek.nl. 4R9ADXsw8evXviF/EWuXbWykotQwJ4F/xHuWYX+p0bm7vGZJQ/rxTp3IGqDcMgpTva0wyorPEHPUViAhugebCg=="),
			test.RRSIG("miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg=="),
			test.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400"),
		},
	}

	m := tc.Msg()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fm.ServeDNS(ctx, rec, m)
	}
}

const dbMiekNLSigned = `
miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1655818598 14400 3600 604800 14400
miek.nl.	1800	IN	RRSIG	SOA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4vFVyKPhIj4RcFeQgbXN3gosojr5iWRfNfjDEUvOfuF3mNQ8sJa1r2eooQ/vV7GLNSFGI03ZtJPGiBetY/5pJg==
miek.nl.	1800	IN	NS	ext.ns.whyscream.net.
miek.nl.	1800	IN	NS	omval.tednet.nl.
miek.nl.	1800	IN	NS	linode.atoom.net.
miek.nl.	1800	IN	NS	ns-ext.nlnetlabs.nl.
miek.nl.	1800	IN	RRSIG	NS 13 2 1800 20220723181226 20220621032053 59725 miek.nl. Wy/bvt7X8qZ1ZOQsHJk+K568VNiv30XvG01h5jq2lvOwh4xc+fivLlOhi8rAV0jduuAUZdu50SSMlQ3XrXYE/Q==
miek.nl.	1800	IN	AAAA	2a01:7e00::f03c:91ff:fef1:6735
miek.nl.	1800	IN	DNSKEY	257 3 13 sfzRg5nDVxbeUc51su4MzjgwpOpUwnuu81SlRHqJuXe3SOYOeypR69tZ52XLmE56TAmPHsiB8Rgk+NTpf0o1Cw==
miek.nl.	1800	IN	CDS	59725 13 1 F8B20281EBD6CDE32FE6C2CAC8CADBAFBE8A3107
miek.nl.	1800	IN	CDS	59725 13 2 57C5D78019685A6575BD3AF29F2156CAEAEF9713DD0F45C357C897242E19BCAF
miek.nl.	1800	IN	CDNSKEY	257 3 13 sfzRg5nDVxbeUc51su4MzjgwpOpUwnuu81SlRHqJuXe3SOYOeypR69tZ52XLmE56TAmPHsiB8Rgk+NTpf0o1Cw==
miek.nl.	14400	IN	NSEC	a.miek.nl. A NS SOA MX AAAA RRSIG NSEC DNSKEY CDS CDNSKEY
miek.nl.	1800	IN	RRSIG	MX 13 2 1800 20220723181226 20220621032053 59725 miek.nl. m21ROZl3a4oiQCP2nf/HmWvrCv99ynVw1ngJvMSOCj80X1zAV64c4v5ZC9J4ARF6l3T2KlCzsgior8vG1vSJdA==
miek.nl.	1800	IN	RRSIG	AAAA 13 2 1800 20220723181226 20220621032053 59725 miek.nl. oE5BbhEnBdGxWXKqt102RAce/Ef/PllZXXu2jPImKw1ouAAdcAXa4uGJAZVfhPSHM4PvtoGIxifgb0bXBc6TJw==
miek.nl.	1800	IN	RRSIG	DNSKEY 13 2 1800 20220723181226 20220621032053 59725 miek.nl. W6SGzSE07xN7lyK7+TTeAGzSJqSPdcvLA1kacA4Vzxz9BAQWwHhiofCYk/27h5PDp5GiYnu3P3frBJr3djo6BA==
miek.nl.	1800	IN	RRSIG	CDS 13 2 1800 20220723181226 20220621032053 59725 miek.nl. rjtNwXIGXevvnrmJ9i4YGwrpB2yBw9WwgQFSSQ63DgGkuvLPFWtMDD4er0R1rGwgXxmIH2u/lXfIsTcov+h+7A==
miek.nl.	1800	IN	RRSIG	CDNSKEY 13 2 1800 20220723181226 20220621032053 59725 miek.nl. w/GP6BMfcUFehd7qK4rcQPeKpTeSt66tXKqHkvw0PcyRfkfMeP0wAjfhj3JXlvhxtPzQTm8nDecXxYwLJEzSUw==
miek.nl.	14400	IN	RRSIG	NSEC 13 2 14400 20220723181226 20220621032053 59725 miek.nl. 4R9ADXsw8evXviF/EWuXbWykotQwJ4F/xHuWYX+p0bm7vGZJQ/rxTp3IGqDcMgpTva0wyorPEHPUViAhugebCg==
miek.nl.	1800	IN	RRSIG	A 13 2 1800 20220723181226 20220621032053 59725 miek.nl. 4+TJBLTRygzRa2/MMWVNFKZbWP+iZVtqvwInKaxVdnjMD/E3zgI0KqnaF/BN7F1j+ib6Hhthri/LpHuGh79LPA==
miek.nl.	1800	IN	A	139.162.196.78
miek.nl.	1800	IN	MX	1 aspmx.l.google.com.
miek.nl.	1800	IN	MX	5 alt1.aspmx.l.google.com.
miek.nl.	1800	IN	MX	5 alt2.aspmx.l.google.com.
miek.nl.	1800	IN	MX	10 aspmx2.googlemail.com.
miek.nl.	1800	IN	MX	10 aspmx3.googlemail.com.
_bf0.miek.nl.	14400	IN	RRSIG	TXT 13 3 14400 20220723181226 20220621032053 59725 miek.nl. zChRSOfjSg+adtj6M1Of7B4seFHVJB+OQ1O7m5ob9cvZhLi+2ADAP6cd7JKDrzOd1LQnyMEK1Vvy+OpI222MfQ==
_bf0.miek.nl.	14400	IN	TXT	"iE0SZVIQp2ZV" "144" "5"
_bf1.miek.nl.	14400	IN	TXT	"arZMWIcKXUNN" "144" "5"
_bf1.miek.nl.	14400	IN	RRSIG	TXT 13 3 14400 20220723181226 20220621032053 59725 miek.nl. QUe5qipncsgjDTrgQ/9DLQCyw0YM6qLQWnesVchJueQMSSbXELDikaqb7naD80SF51Jjpkr8REnUZqfe0o0MXA==
a.miek.nl.	1800	IN	A	139.162.196.78
a.miek.nl.	1800	IN	AAAA	2a01:7e00::f03c:91ff:fef1:6735
a.miek.nl.	14400	IN	NSEC	archive.miek.nl. A AAAA RRSIG NSEC
a.miek.nl.	1800	IN	RRSIG	A 13 3 1800 20220723181226 20220621032053 59725 miek.nl. bJ4mHlMAbl6MdBnhyZavM7fSE3clBLQKXeDuVW0keERG6HtX3VWd3jfylZPfAzVc4Y9zOuSUW87QkEtpiG0Wlw==
a.miek.nl.	1800	IN	RRSIG	AAAA 13 3 1800 20220723181226 20220621032053 59725 miek.nl. dphM/QpdM6R5cgNb0ZG+x27IQD3Ytzbu7/6NqcGfSINywY9GVaA6RY61b0PNcDiA2RCFbwMK7IIysUirLZqNfg==
a.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. aNBNT2VX17Pa716Gi+Fc/H4ZeSrjsCR53U44EJz9EIBYMA1urlhEeGi7g0H9GiyVxKOjbBSc51PWatV+UCUBRA==
archive.miek.nl.	1800	IN	CNAME	a.miek.nl.
archive.miek.nl.	14400	IN	NSEC	go.dns.miek.nl. CNAME RRSIG NSEC
archive.miek.nl.	1800	IN	RRSIG	CNAME 13 3 1800 20220723181226 20220621032053 59725 miek.nl. 16SMEgPqEZ9J8LZjbRDrQ/92XpcNMrwTeMaLk13WUSLVvIMmB+CJNm56hPt35LmsB0DxoGH/+uP4A1teMQ29HQ==
archive.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. Sqdlq6upOFmgFUmCQ1hUNRCX41qwHhENHjTEk2ctLgRoJIGVI6KP3P/viqQYdcesHSxO8xo8NV3F9NZoPSi2JA==
go.dns.miek.nl.	1800	IN	RRSIG	TXT 13 4 1800 20220723181226 20220621032053 59725 miek.nl. qCbh3ODMprlaY6LZ7YMs4X/CR7A7ucn6AqusCx2pxM1b05MYQTF5P5vxAPiTIKLF1n8I53507gAOqO+AOvDpbw==
go.dns.miek.nl.	14400	IN	RRSIG	NSEC 13 4 14400 20220723181226 20220621032053 59725 miek.nl. Ba3IsA0gilxGtIag0SDRdR4EYgw5gaS6x0FsibEMFfLtcbwPiYhqjpzFo1PX1BIKUKv0+vuwXO+sVf+P6gJHig==
go.dns.miek.nl.	1800	IN	TXT	"Hello!"
go.dns.miek.nl.	14400	IN	NSEC	www.miek.nl. TXT RRSIG NSEC
www.miek.nl.	1800	IN	CNAME	a.miek.nl.
www.miek.nl.	14400	IN	NSEC	miek.nl. CNAME RRSIG NSEC
www.miek.nl.	1800	IN	RRSIG	CNAME 13 3 1800 20220723181226 20220621032053 59725 miek.nl. 35vdb9u8Yxl92YpwVxB0NcnMO6xTrJ7eN2cuYppur5Vu74QbxwifHr9zxmLVnORPEVe80mqviXuMweu1GyvcCA==
www.miek.nl.	14400	IN	RRSIG	NSEC 13 3 14400 20220723181226 20220621032053 59725 miek.nl. 1+g3+5nxTZC5CDkC1Mu8hmNW56Rr8YbRCPjjGDqZd71e8SUYNshj+bjSJSncXgm96CZeYFj02uV2ISpIf+VWig==
`
