package file_nsec5

import (
	"crypto/rand"
	"encoding/base64"
)

func max(a, b uint64) uint64 {
	if a <= b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// creates a slice of random bytes of length s (used for testing)
func randomBytes(s int) []byte {
	b := make([]byte, s)
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

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func fromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
