package bloomsec

import (
	"os"
	"testing"

	"github.com/coredns/coredns/plugin/bloomfile"
)

func TestNames(t *testing.T) {
	f, err := os.Open("testdata/db.miek.nl_ns")
	if err != nil {
		t.Error(err)
	}
	z, err := bloomfile.Parse(f, "db.miek.nl_ns", "miek.nl", 0)
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
	chunkSize := 14 * 255 * 8
	b := make([]bool, chunkSize)
	k := uint64(20)
	m := uint64(10)
	for i := 0; i < int(chunkSize); i++ {
		b[i] = randBool()
	}
	calc, r1, r2, err := stringsToBits(bitsToStrings(&b, 14, m, k), uint64(chunkSize))
	if err != nil {
		t.Error("We got an error")
	}
	if m != r1 {
		t.Errorf("m not equal")
	}
	if k != r2 {
		t.Errorf("k not equal")
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
