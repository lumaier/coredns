package bloomfile_nsec5

import (
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"strconv"

	"github.com/coredns/coredns/plugin/bloomfile_nsec5/rrutil"
)

type bfChunk struct {
	globalIndex uint64
	bitArray    []bool
}

type bloomfilter struct {
	m                  uint64 // size of bitarray
	k                  uint64 // number of hashfunctions used
	l                  uint64 // number of invocations used of SHA512
	bitArray           []bool // Bloomfilter
	nr_bytes_per_index uint64 // number of bytes of SHA512 output used to create index
	sha512hasher       hash.Hash
}

func newBloomfilter(m, k uint64) *bloomfilter {
	bf := bloomfilter{}
	bf.m = m
	bf.k = k
	bf.sha512hasher = sha512.New()

	// we use at least 4 bytes s.t. we can create an uint32 of it
	log_m := math.Ceil(math.Log2(float64(bf.m)))
	bf.nr_bytes_per_index = rrutil.Max(4, rrutil.NextPowOf2(uint64(math.Ceil(log_m/8.))))

	// we calculate the number of invocations we need of SHA512
	// each invocation gives 64 Bytes
	bf.l = uint64(math.Ceil(float64(bf.k*bf.nr_bytes_per_index) / 64.))

	bf.bitArray = make([]bool, bf.m)

	return &bf
}

// returns false if at least one of the corresponding bits is 0 and additionally returns also this index
// returns 0 as index if all bits were 1
func (bf *bloomfilter) lookup(element []byte) (uint64, bool) {
	indices := bf.calculateIndices(element)
	for _, index := range indices {
		if !bf.bitArray[index] {
			return index, false
		}
	}
	return 0, true
}

// returns the array of k indices of type uint64
func (bf *bloomfilter) calculateIndices(data []byte) []uint64 {
	digest := []byte{}
	for i := uint64(1); i <= bf.l; i++ {
		x := append(data, byte(i))
		bf.sha512hasher.Write(x)
		digest = bf.sha512hasher.Sum(digest)
	}

	indices := make([]uint64, bf.k)
	for i := uint64(0); i < bf.k; i++ {
		if bf.nr_bytes_per_index >= 8 {
			// FIXME: if bf.nr_bytes_per_index is > 8 then all bytes with index >= 8 are ignored. is this legal?
			indices[i] = binary.BigEndian.Uint64(digest[i*bf.nr_bytes_per_index:(i+1)*bf.nr_bytes_per_index]) % bf.m
		} else { // then bf.nr_byte_per_index is guaranteed to be 4 (power of 2 and at least 4)
			indices[i] = uint64(binary.BigEndian.Uint32(digest[i*bf.nr_bytes_per_index:(i+1)*bf.nr_bytes_per_index]) % uint32(bf.m))
		}

	}

	return indices
}

func isBfTxtChunk(origin, name string) bool {
	k := len(name)
	n := len(origin)
	if name[:3] == "_bf" && name[k-n:] == origin {
		return true
	}
	return false
}

func extractGlobalIndex(origin, name string) (uint64, error) {
	k := len(name)
	n := len(origin)
	i := name[3 : k-n-1]
	r, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return r, nil
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

func (bf *bloomfilter) PrintWhole() string {
	result := "\n["
	for _, x := range bf.bitArray {
		if x {
			result = result + "1 "
		} else {
			result = result + "0 "
		}
	}
	return result + "]\n\n"
}

func (bf *bloomfilter) Info() string {
	result := fmt.Sprintf("Bloomfilter capacity: %d (bits needed to index: %d)\n", bf.m, int(math.Ceil(math.Log2(float64(bf.m)))))
	result += fmt.Sprintf("number of invocations of SHA512: %d\n", bf.l)
	result += fmt.Sprintf("Bloomfilter occupancy: %f\n", bf.Occupancy())
	return result
}

func (bf *bloomfilter) Occupancy() float64 {
	n := uint64(0)
	for _, b := range bf.bitArray {
		if b {
			n++
		}
	}
	return float64(n) / float64(bf.m)
}

func printIndices(indices []uint64) string {
	result := ""
	for _, i := range indices {
		result += fmt.Sprintf("%d ", i)
	}
	return result + "\n"
}
