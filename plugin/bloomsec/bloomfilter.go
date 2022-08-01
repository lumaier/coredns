package bloomsec

import (
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
)

type bfChunk struct {
	globalIndex uint64
	bitArray    []bool
	m           uint64
	k           uint64
}

type bloomfilter struct {
	n                  uint64  // number of items in the Bloomfilter
	m                  uint64  // size of bitarray
	k                  uint64  // number of indices used
	l                  uint64  // number of invocations used of SHA512
	bitArray           []bool  // Bloomfilter
	fprate             float64 // actual estimated false positive rate (given n, m and k)
	nr_bytes_per_index uint64  // number of bytes of SHA512 output used to create index
	chunkSize          uint64  // number of bits in bitarray per chunk
	sha512hasher       hash.Hash
}

// creates a new bloom filter
func newBloomfilter(n uint64, falsePositiveProb float64, chunkSize uint64) *bloomfilter {
	bf := bloomfilter{}
	bf.n = n
	bf.chunkSize = chunkSize
	bf.sha512hasher = sha512.New()
	bf.estimateParameters(falsePositiveProb)
	// set m to the next multiple of chunkSize (for nice chunking)
	bf.m = uint64(math.Ceil((float64(bf.m) / float64(chunkSize)))) * chunkSize

	// (1-e**(-kn/m))**k
	bf.fprate = math.Pow(1-math.Exp(-float64(bf.k*bf.n)/float64(bf.m)), float64(bf.k))

	// we use at least 4 bytes s.t. we can create an uint32 of it
	log_m := math.Ceil(math.Log2(float64(bf.m)))
	bf.nr_bytes_per_index = max(4, nextPowOf2(uint64(math.Ceil(log_m/8.))))

	// we calculate the number of invocations we need of SHA512
	// each invocation gives 64 Bytes
	bf.l = uint64(math.Ceil(float64(bf.k*bf.nr_bytes_per_index) / 64.))

	bf.bitArray = make([]bool, bf.m)

	return &bf
}

// sets all corresponding bits to 1
func (bf *bloomfilter) insert(element []byte) {
	indices := bf.calculateIndices(element)
	for _, index := range indices {
		bf.bitArray[index] = true
	}
}

// returns false if at least one of the corresponding bits is 0
func (bf *bloomfilter) lookup(element []byte) bool {
	indices := bf.calculateIndices(element)
	for _, index := range indices {
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

// returns the array of k indices of type uint64
func (bf *bloomfilter) calculateIndices(data []byte) []uint64 {
	digest := []byte{}
	for i := uint64(1); i <= bf.l; i++ {
		x := append(data, byte(i))
		//bf.sha512hasher.Write(x)
		//digest = bf.sha512hasher.Sum(digest)
		var temp = sha512.Sum512(x)
		digest = append(digest, temp[:]...)
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

func (bf *bloomfilter) estimateParameters(fp float64) {
	nf := float64(bf.n)
	ln2 := math.Log(2)
	num := nf * math.Log(fp)
	deno := math.Pow(ln2, 2)

	//m = -1 * (n * lnP)/(ln2)^2
	bf.m = uint64(-1 * num / deno)
	// k = m/n * ln2
	bf.k = uint64(math.Ceil(float64(bf.m) / nf * ln2))
}

func (bf *bloomfilter) print() {
	fmt.Printf("\nNumber of insertable items: %d\n", bf.n)
	fmt.Printf("Bloomfilter capacity: %d (bits needed to index: %d)\n", bf.m, int(math.Ceil(math.Log2(float64(bf.m)))))
	fmt.Printf("number of invocations of SHA512: %d\n", bf.l)
	fmt.Printf("Bloomfilter false positiverate: %f\n", bf.fprate)
	fmt.Printf("Bloomfilter occupancy: %f\n\n", bf.Occupancy())
}

func (bf *bloomfilter) Info() string {
	result := fmt.Sprintf("\nNumber of insertable items: %d\n", bf.n)
	result += fmt.Sprintf("Bloomfilter capacity: %d (bits needed to index: %d)\n", bf.m, int(math.Ceil(math.Log2(float64(bf.m)))))
	result += fmt.Sprintf("number of invocations of SHA512: %d\n", bf.l)
	result += fmt.Sprintf("Bloomfilter false positiverate: %f\n", bf.fprate)
	result += fmt.Sprintf("Bloomfilter occupancy: %f\n\n", bf.Occupancy())
	return result
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

func (bf *bloomfilter) Occupancy() float64 {
	n := uint64(0)
	for _, b := range bf.bitArray {
		if b {
			n++
		}
	}
	return float64(n) / float64(bf.m)
}

// returns the false positive rate when n elements are stored,
// runs 100000 tests,
// ALARM: this clears the Bloomfilter
func (bf *bloomfilter) falsePositiveRate(n uint64) float64 {
	nr_tests := uint64(100000)
	bf.bitArray = make([]bool, bf.m)
	i := uint64(0)
	b := make([]byte, 8)
	for ; i < n; i++ {
		binary.BigEndian.PutUint64(b, i)
		bf.insert(b)
	}
	fp := uint64(0)
	for ; i < n+nr_tests; i++ {
		binary.BigEndian.PutUint64(b, i)
		if bf.lookup(b) {
			fp++
		}
	}
	return float64(fp) / float64(nr_tests)
}

// Divides the Bloomfilter into equally sized chunks (of length bf.chunkSize) and returns them.
// The global index corresponds to the position in the global Bloom filter (indexed 0, 1, 2, etc...)
func (bf *bloomfilter) chunking() (*[]bfChunk, error) {
	n_chunks := bf.m / bf.chunkSize
	if bf.m%bf.chunkSize != 0 {
		return nil, fmt.Errorf("Bloom filter could not be chunked into equally sized chunks.")
	}
	chunks := make([]bfChunk, n_chunks)
	for i := uint64(0); i < n_chunks; i++ {
		chunks[i] = bfChunk{i, bf.bitArray[i*bf.chunkSize : (i+1)*bf.chunkSize], bf.m, bf.k}
	}
	return &chunks, nil
}
