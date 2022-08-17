package file_nsec5

import (
	"testing"
	"time"

	"github.com/ProtonMail/go-ecvrf/ecvrf"
)

const N = 200

const vrf_pubkey = "L+9spXqxohFCjVQNAlEZziifcnvdDjmbZO/Hl86pCo4="
const vrf_privkey = "XFKVYiagDkqCbNTXmTxedPOGD87Prc3IvuPebYEcyFAv72ylerGiEUKNVA0CURnOKJ9ye90OOZtk78eXzqkKjg=="

func TestResolver(t *testing.T) {
	s, _ := fromBase64(vrf_pubkey)
	pubkey, _ := ecvrf.NewPublicKey(s)
	output := make([][][]byte, N*2)
	for i, _ := range output {
		output[i] = create_output()
	}
	start := time.Now()
	for i := 0; i < N; i++ {
		for j := 0; j < 2; j++ {
			_, vrf, _ := pubkey.Verify(output[2*i+j][0], output[2*i+j][2])
			Equal(output[i][2], vrf)
		}
	}
	log.Infof("%d queries elapsed in %d mikrosec", N, time.Since(start).Nanoseconds()/1000)
}

func create_output() []([]byte) {
	s, _ := fromBase64(vrf_privkey)
	privkey, _ := ecvrf.NewPrivateKey(s)
	t := randomBytes(10)
	vrf, proof, _ := privkey.Prove(t)
	return []([]byte){t, vrf, proof}
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
