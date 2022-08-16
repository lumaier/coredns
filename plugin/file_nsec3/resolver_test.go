package file_nsec3

import (
	"crypto/sha1"
	"testing"
	"time"
)

const N = 10000

func TestResolver(t *testing.T) {
	output := make([][][]byte, N)
	for i, _ := range output {
		output[i] = create_output()
	}
	start := time.Now()
	for i := 0; i < N; i++ {
		for j := 0; j < 2; j++ {
			hash := sha1.Sum(output[i][0])
			Equal(output[i][1], hash[:])
		}
	}
	log.Infof("%d queries elapsed in %d mikrosec", N, time.Since(start).Nanoseconds()/1000)
}

func create_output() []([]byte) {
	t := randomBytes(10)
	temp := sha1.Sum(t)
	return []([]byte){t, temp[:]}
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
