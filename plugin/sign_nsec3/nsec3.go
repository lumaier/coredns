package sign_nsec3

import (
	"github.com/miekg/dns"
)

func NSEC3(domain, name, next string, ttl uint32) *dns.TXT {

	return &dns.TXT{
		Hdr: dns.RR_Header{Name: domain, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: []string{"nsec3", name, next},
	}
}
