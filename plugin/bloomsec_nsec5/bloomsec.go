package bloomsec_nsec5

import (
	"fmt"
	"sort"

	"github.com/coredns/coredns/plugin/bloomfile_nsec5"
	"github.com/coredns/coredns/plugin/bloomfile_nsec5/tree"

	"github.com/miekg/dns"
)

const (
	bitsEncoded = 24 // results in 4 chars
)

// names returns the elements of the zone in nsec order.
func names(z *bloomfile_nsec5.Zone) []string {
	// There will also be apex records other than NS and SOA (who are kept separate), as we
	// are adding DNSKEY and CDS/CDNSKEY records in the apex *before* we sign.
	n := []string{}
	z.AuthWalk(func(e *tree.Elem, _ map[uint16][]dns.RR, auth bool) error {
		if !auth {
			return nil
		}
		n = append(n, e.Name())
		return nil
	})
	return n
}

// Given a Bloom filter chunk it returns the corresponding TXT record.
// It assumes that the bitarray length of the chunk is a multiple of 24
func bloomTXT(apexname string, chunk *bfChunk, ttl uint32) (*dns.TXT, error) {
	name := "_bf" + fmt.Sprint(chunk.globalIndex) + "." + apexname

	if len(chunk.bitArray)%24 != 0 {
		return nil, fmt.Errorf("The bitarray length in the chunk %d is not a multiple of 24.\n", chunk.globalIndex)
	}

	bitarray_bytes, err := bitsToBytes(&chunk.bitArray)
	if err != nil {
		return nil, err
	}
	encoded_bitarray := toBase64(*bitarray_bytes)

	txt_data := []string{}
	for i := 0; i < len(encoded_bitarray); i += 255 {
		txt_data = append(txt_data, encoded_bitarray[i:min(i+255, len(encoded_bitarray))])
	}

	txt_data = append(txt_data, fmt.Sprint(chunk.m), fmt.Sprint(chunk.k))

	return &dns.TXT{
		Hdr: dns.RR_Header{Name: name, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: txt_data,
	}, nil
}

// Given a slice of bits it returns a slice of bytes encoding the slice of bits. len(b) must be a multiple of 8.
func bitsToBytes(b *[]bool) (*[]byte, error) {
	l := len(*b)
	if l%8 != 0 {
		return nil, fmt.Errorf("The bit slice length is not a multiple of 8.\n")
	}
	l /= 8
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		for j := 0; j < 8; j++ {
			if (*b)[i*8+j] {
				bytes[i] |= 0x80 >> uint(j)
			}
		}
	}

	return &bytes, nil
}

// Decodes a slice of bytes into a slice of bits
func bytesToBits(bytes *[]byte) *[]bool {
	l := len(*bytes)
	bits := make([]bool, l*8)
	for i := 0; i < l; i++ {
		for j := 0; j < 8; j++ {
			if ((*bytes)[i]<<uint(j))&0x80 == 0x80 {
				bits[i*8+j] = true
			}
		}
	}
	return &bits
}

// Returns an NSEC record according to name, next, ttl and bitmap. Note that the bitmap is sorted before use.
func NSEC(name, next string, ttl uint32, bitmap []uint16) *dns.NSEC {
	sort.Slice(bitmap, func(i, j int) bool { return bitmap[i] < bitmap[j] })

	return &dns.NSEC{
		Hdr:        dns.RR_Header{Name: name, Ttl: ttl, Rrtype: dns.TypeNSEC, Class: dns.ClassINET},
		NextDomain: next,
		TypeBitMap: bitmap,
	}
}
