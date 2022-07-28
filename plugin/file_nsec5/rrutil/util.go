// Package rrutil provides function to find certain RRs in slices.
package rrutil

import (
	"crypto/rand"
	"encoding/base64"
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

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func FromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
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
