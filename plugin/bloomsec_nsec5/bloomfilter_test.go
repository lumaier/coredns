package bloomsec_nsec5

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

func TestChunking(t *testing.T) {
	n := 10000
	fp := 0.001
	chunkSize := uint64(50)
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
	chunks, err := bf.chunking()
	if err != nil {
		t.Error("Error")
	}
	if uint64(len(*chunks)) != bf.m/chunkSize {
		t.Errorf("not the correct number of chunks")
	}
	for _, c := range *chunks {
		if len(c.bitArray) != int(chunkSize) {
			t.Errorf("not the correct chunk length")
		}
	}

}
