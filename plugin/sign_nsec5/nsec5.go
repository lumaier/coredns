package sign_nsec5

import (
	"github.com/miekg/dns"
)

const (
	bitsEncoded = 24 // results in 4 chars
)

type VRF_output struct {
	hash   string
	proof  string
	domain string
}

type NSEC5RR struct {
	nsec5      *dns.TXT
	nsec5proof *dns.TXT
}

// Returns a TXT record containing an NSEC5 record and a TXT record containing the NSEC5PROOF
// wildcard bit is not implemented
func NSEC5(domain, name, proof, next string, ttl uint32) (*dns.TXT, *dns.TXT) {

	r1 := &dns.TXT{
		Hdr: dns.RR_Header{Name: domain, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: []string{"nsec5", name, next},
	}

	r2 := &dns.TXT{
		Hdr: dns.RR_Header{Name: domain, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: []string{"nsec5proof", proof},
	}

	return r1, r2
}
