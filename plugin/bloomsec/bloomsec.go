package bloomsec

import (
	"fmt"

	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/file/tree"

	"github.com/miekg/dns"
)

// names returns the elements of the zone in nsec order.
func names(z *file.Zone) []string {
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

// bloomTXT: given a bloomfilter chunk it returns the corresponding TXT record
// it assumes that the bitarray length is equal the chunkSize
func bloomTXT(apexname string, chunk *bfChunk, ttl uint32) (*dns.TXT, error) {
	name := "_bf" + fmt.Sprint(chunk.globalIndex) + "." + apexname

	n_strings := len(chunk.bitArray) / (255 * 8)
	if n_strings != (int(chunkSize) / (255 * 8)) {
		return nil, fmt.Errorf("")
	}

	strings := bitsToStrings(&chunk.bitArray, n_strings)

	return &dns.TXT{
		Hdr: dns.RR_Header{Name: name, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: *strings,
	}, nil
}

func bitsToStrings(b *[]bool, n_strings int) *[]string {
	strings := make([]string, n_strings)
	for i := 0; i < n_strings; i++ {
		temp := make([]byte, 255)
		for j := 0; j < 255; j++ {
			for k := 0; k < 8; k++ {
				if (*b)[i*(255*8)+j*8+k] {
					temp[j] |= 0x80 >> uint(k%8)
				}
			}
		}
		strings[i] = string(temp)
	}
	return &strings
}

func stringsToBits(strings *[]string) *[]bool {
	b := make([]bool, chunkSize)
	var temp byte
	for i := 0; i < len(*strings); i++ {
		for j := 0; j < 255; j++ {
			temp = byte((*strings)[i][j])
			for k := 0; k < 8; k++ {
				if (temp<<uint(k%8))&0x80 == 0x80 {
					b[i*(255*8)+j*8+k] = true
				}
			}
		}
	}
	return &b
}
