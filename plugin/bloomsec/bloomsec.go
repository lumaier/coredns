package bloomsec

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/coredns/coredns/plugin/bloomfile"
	"github.com/coredns/coredns/plugin/bloomfile/tree"

	"github.com/miekg/dns"
)

// names returns the elements of the zone in nsec order.
func names(z *bloomfile.Zone) []string {
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
// It assumes that the bitarray length is equal the chunkSize
func bloomTXT(apexname string, chunk *bfChunk, ttl uint32) (*dns.TXT, error) {
	name := "_bf" + fmt.Sprint(chunk.globalIndex) + "." + apexname

	n_strings := len(chunk.bitArray) / (255 * 8)
	if n_strings != (int(chunkSize) / (255 * 8)) {
		return nil, fmt.Errorf("The bitarray in the chunk %d has not the correct length.", chunk.globalIndex)
	}

	strings := bitsToStrings(&chunk.bitArray, n_strings, chunk.m, chunk.k)

	return &dns.TXT{
		Hdr: dns.RR_Header{Name: name, Ttl: ttl, Rrtype: dns.TypeTXT, Class: dns.ClassINET},
		Txt: *strings,
	}, nil
}

// Given a slice of bits it returns a slice of strings encoding the slice of bits. The last two strings contain the length of the
// Bloom filter (m) and the number of indices used (k)
func bitsToStrings(b *[]bool, n_strings int, m, k uint64) *[]string {
	strings := make([]string, n_strings+2)
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
	strings[n_strings] = fmt.Sprint(m)
	strings[n_strings+1] = fmt.Sprint(k)
	return &strings
}

// Decodes a slice of strings into a slice of bits, the total length of the Bloom filter and the number of indices used.
// Throws an error when the last two strings can't be parsed into an uint64.
func stringsToBits(strings *[]string) (*[]bool, uint64, uint64, error) {
	b := make([]bool, chunkSize)
	var temp byte
	l := len(*strings) - 2 // the last two strings correspond to m and k
	for i := 0; i < l; i++ {
		for j := 0; j < 255; j++ {
			temp = byte((*strings)[i][j])
			for k := 0; k < 8; k++ {
				if (temp<<uint(k%8))&0x80 == 0x80 {
					b[i*(255*8)+j*8+k] = true
				}
			}
		}
	}
	result1, err := strconv.ParseUint((*strings)[l], 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	result2, err := strconv.ParseUint((*strings)[l+1], 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	return &b, result1, result2, err
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
