package bloomsec

import (
	"crypto/rand"
)

func max(a, b uint64) uint64 {
	if a <= b {
		return b
	}
	return a
}

func randomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func nextPowOf2(n uint64) uint64 {
	k := uint64(1)
	for k < n {
		k = k << 1
	}
	return uint64(k)
}

func randBool() bool {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]&0x80 == 0x80
}
