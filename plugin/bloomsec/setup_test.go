package bloomsec

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		exp       *Signer
	}{
		{`bloomsec testdata/db.miek.nl miek.nl {
			key file testdata/Kmiek.nl.+013+59725
			fp_rate 0.02
		 }`,
			false,
			&Signer{
				keys:       []Pair{},
				origin:     "miek.nl.",
				dbfile:     "testdata/db.miek.nl",
				directory:  "/var/lib/coredns",
				signedfile: "db.miek.nl.signed",
				fp_prob:    0.02,
				chunkSize:  14 * 255 * 8,
			},
		},
		{`bloomsec testdata/db.miek.nl example.org {
			key file testdata/Kmiek.nl.+013+59725
			directory testdata
			chunksize 2040
		 }`,
			false,
			&Signer{
				keys:       []Pair{},
				origin:     "example.org.",
				dbfile:     "testdata/db.miek.nl",
				directory:  "testdata",
				signedfile: "db.example.org.signed",
				chunkSize:  2040,
				fp_prob:    0.001,
			},
		},
		// errors
		{`bloomsec db.example.org {
			key file /etc/coredns/keys/Kexample.org
		 }`,
			true,
			nil,
		},
	}
	for i, tc := range tests {
		c := caddy.NewTestController("dns", tc.input)
		sign, err := parse(c)

		if err == nil && tc.shouldErr {
			t.Fatalf("Test %d expected errors, but got no error", i)
		}
		if err != nil && !tc.shouldErr {
			t.Fatalf("Test %d expected no errors, but got '%v'", i, err)
		}
		if tc.shouldErr {
			continue
		}
		signer := sign.signers[0]
		if x := signer.origin; x != tc.exp.origin {
			t.Errorf("Test %d expected %s as origin, got %s", i, tc.exp.origin, x)
		}
		if x := signer.dbfile; x != tc.exp.dbfile {
			t.Errorf("Test %d expected %s as dbfile, got %s", i, tc.exp.dbfile, x)
		}
		if x := signer.directory; x != tc.exp.directory {
			t.Errorf("Test %d expected %s as directory, got %s", i, tc.exp.directory, x)
		}
		if x := signer.signedfile; x != tc.exp.signedfile {
			t.Errorf("Test %d expected %s as signedfile, got %s", i, tc.exp.signedfile, x)
		}
		if x := signer.chunkSize; x != tc.exp.chunkSize {
			t.Errorf("Test %d expected %d as chunkSize, got %d", i, tc.exp.chunkSize, x)
		}
		if x := signer.fp_prob; x != tc.exp.fp_prob {
			t.Errorf("Test %d expected %f as fp_prob, got %f", i, tc.exp.fp_prob, x)
		}
	}
}
