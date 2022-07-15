package bloomsec_nsec5

import (
	"os"
	"testing"

	"github.com/coredns/coredns/plugin/bloomfile_nsec5"
)

func TestNames(t *testing.T) {
	f, err := os.Open("testdata/db.miek.nl_ns")
	if err != nil {
		t.Error(err)
	}
	z, err := bloomfile_nsec5.Parse(f, "db.miek.nl_ns", "miek.nl", 0)
	if err != nil {
		t.Error(err)
	}

	names := names(z)
	expected := []string{"miek.nl.", "child.miek.nl.", "www.miek.nl."}
	for i := range names {
		if names[i] != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], names[i])
		}
	}
}

func TestConversion(t *testing.T) {
	chunkSize := 24 * 163
	b := make([]bool, chunkSize)
	for i := 0; i < int(chunkSize); i++ {
		b[i] = randBool()
	}
	temp, err := bitsToBytes(&b)
	calc := bytesToBits(temp)
	if err != nil {
		t.Error("error")
	}

	if len(*calc) != len(b) {
		t.Errorf("Length not equal")
	}
	for i, x := range b {
		if x != (*calc)[i] {
			t.Errorf("Bit %d not equal", i)
		}
	}
}
