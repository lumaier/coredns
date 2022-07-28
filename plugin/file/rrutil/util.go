// Package rrutil provides function to find certain RRs in slices.
package rrutil

import (
	"fmt"
	"runtime"

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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
