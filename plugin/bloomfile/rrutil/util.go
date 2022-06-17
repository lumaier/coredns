// Package rrutil provides function to find certain RRs in slices.
package rrutil

import (
	"crypto/rand"

	"github.com/miekg/dns"
)

// SubTypeSignature returns the RRSIG for the subtype.
func SubTypeSignature(rrs []dns.RR, subtype uint16) []dns.RR {
	sigs := []dns.RR{}
	// there may be multiple keys that have signed this subtype
	for _, sig := range rrs {
		if s, ok := sig.(*dns.RRSIG); ok {
			if s.TypeCovered == subtype {
				sigs = append(sigs, s)
			}
		}
	}
	return sigs
}

func Max(a, b uint64) uint64 {
	if a <= b {
		return b
	}
	return a
}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func NextPowOf2(n uint64) uint64 {
	k := uint64(1)
	for k < n {
		k = k << 1
	}
	return uint64(k)
}

func RandBool() bool {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]&0x80 == 0x80
}
