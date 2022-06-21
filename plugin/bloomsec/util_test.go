package bloomsec

import "testing"

func TestEncodeDecode(t *testing.T) {
	temp := []byte("ajlkjaölksdjlkajblkasjdlfkkjalskdjflksjflkajölskdjblkakasdfasjsdklfjsl")
	r, err := fromBase64(toBase64(temp))
	if err != nil {
		t.Fatal("error")
	}
	if len(r) != len(temp) {
		t.Fatal("not the same length")
	}
	for i, v := range temp {
		if v != r[i] {
			t.Fatal("slices not identical")
		}
	}
}
