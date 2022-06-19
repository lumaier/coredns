package bloomsec

import (
	"testing"
)

func TestBloomfilter(t *testing.T) {
	n := 1_000_000
	fp := 0.001
	chunkSize := uint64(14 * 255 * 8)
	domainNames := make([][]byte, n)
	for i := 0; i < n; i++ {
		domainNames[i] = randomBytes(32)
	}
	bf := newBloomfilter(uint64(n), fp, chunkSize)
	for _, x := range domainNames {
		bf.insert(x)
		test := bf.lookup(x)
		if !test {
			t.Errorf("FATAL. Should have been in the Bloomfilter")
		}
	}
}
