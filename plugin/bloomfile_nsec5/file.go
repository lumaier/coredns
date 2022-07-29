// Package file implements a file backend.
package bloomfile_nsec5

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/transfer"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/plugin/bloomfile_nsec5/rrutil"
)

var log = clog.NewWithPlugin("bloomfile_nsec5")

type (
	// File is the plugin that reads zone data from disk.
	File struct {
		Next plugin.Handler
		Zones
		transfer *transfer.Transfer
	}

	// Zones maps zone names to a *Zone.
	Zones struct {
		Z     map[string]*Zone // A map mapping zone (origin) to the Zone's data
		Names []string         // All the keys from the map Z as a string slice.
	}
)

// ServeDNS implements the plugin.Handle interface.
func (f File) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	qname := state.Name()
	// TODO(miek): match the qname better in the map
	zone := plugin.Zones(f.Zones.Names).Matches(qname)
	if zone == "" {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}

	z, ok := f.Zones.Z[zone]
	if !ok || z == nil {
		return dns.RcodeServerFailure, nil
	}

	// If transfer is not loaded, we'll see these, answer with refused (no transfer allowed).
	if state.QType() == dns.TypeAXFR || state.QType() == dns.TypeIXFR {
		return dns.RcodeRefused, nil
	}

	// This is only for when we are a secondary zones.
	if r.Opcode == dns.OpcodeNotify {
		if z.isNotify(state) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Authoritative = true
			w.WriteMsg(m)

			log.Infof("Notify from %s for %s: checking transfer", state.IP(), zone)
			ok, err := z.shouldTransfer()
			if ok {
				z.TransferIn()
			} else {
				log.Infof("Notify from %s for %s: no SOA serial increase seen", state.IP(), zone)
			}
			if err != nil {
				log.Warningf("Notify from %s for %s: failed primary check: %s", state.IP(), zone, err)
			}
			return dns.RcodeSuccess, nil
		}
		log.Infof("Dropping notify from %s for %s", state.IP(), zone)
		return dns.RcodeSuccess, nil
	}

	z.RLock()
	exp := z.Expired
	z.RUnlock()
	if exp {
		log.Errorf("Zone %s is expired", zone)
		return dns.RcodeServerFailure, nil
	}

	answer, ns, extra, result := z.Lookup(ctx, state, qname)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer, m.Ns, m.Extra = answer, ns, extra

	switch result {
	case Success:
	case NoData:
	case NameError:
		m.Rcode = dns.RcodeNameError
	case Delegation:
		m.Authoritative = false
	case ServerFailure:
		// If the result is SERVFAIL and the answer is non-empty, then the SERVFAIL came from an
		// external CNAME lookup and the answer contains the CNAME with no target record. We should
		// write the CNAME record to the client instead of sending an empty SERVFAIL response.
		if len(m.Answer) == 0 {
			return dns.RcodeServerFailure, nil
		}
		//  The rcode in the response should be the rcode received from the target lookup. RFC 6604 section 3
		m.Rcode = dns.RcodeServerFailure
	}

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (f File) Name() string { return "file" }

type serialErr struct {
	err    string
	zone   string
	origin string
	serial int64
}

func (s *serialErr) Error() string {
	return fmt.Sprintf("%s for origin %s in file %s, with %d SOA serial", s.err, s.origin, s.zone, s.serial)
}

// Parse parses the zone in filename and returns a new Zone or an error.
// If serial >= 0 it will reload the zone, if the SOA hasn't changed
// it returns an error indicating nothing was read.
func Parse(f io.Reader, origin, fileName, vrfFileName string, serial int64) (*Zone, error) {
	zp := dns.NewZoneParser(f, dns.Fqdn(origin), fileName)
	zp.SetIncludeAllowed(true)
	z := NewZone(origin, fileName)

	if vrfFileName != "nil" {
		z.ReadKeys(vrfFileName)
	}

	seenSOA := false

	// this list holds all chunks that correspond to bloom filter chunks in TXT records including its global index
	chunks := [](*dns.TXT){}

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		if err := zp.Err(); err != nil {
			return nil, err
		}

		if !seenSOA {
			if s, ok := rr.(*dns.SOA); ok {
				seenSOA = true

				// -1 is valid serial is we failed to load the file on startup.

				if serial >= 0 && s.Serial == uint32(serial) { // same serial
					return nil, &serialErr{err: "no change in SOA serial", origin: origin, zone: fileName, serial: serial}
				}
			}
		}

		if err := z.Insert(rr); err != nil {
			return nil, err
		}

		// check whether it is a bloomfilter chunk, an nsec5 or nsec5proof
		if s, ok := rr.(*dns.TXT); ok && isBfTxtChunk(origin, rr.Header().Name) {
			chunks = append(chunks, s)
		} else if ok && (*s).Txt[0] == "nsec5" {
			z.nsec5s = append(z.nsec5s, s)
		}
	}

	sort.Slice(z.nsec5s, func(i, j int) bool {
		return z.nsec5s[i].Txt[1] < z.nsec5s[j].Txt[1]
	})

	z.N_nsec5s = len(z.nsec5s)

	if !seenSOA {
		return nil, fmt.Errorf("file %q has no SOA record for origin %s", fileName, origin)
	}

	hasBF := len(chunks) > 0
	var m, k uint64
	var err error
	var b *[]bool

	if hasBF {
		b, m, k, err = decodeBloomTxt(chunks[0])
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		z.bf = *newBloomfilter(m, k)
		z.chunkSize = uint64(len(*b))
	}

	if len(chunks)*int(z.chunkSize) != int(z.bf.m) {
		return nil, fmt.Errorf("not all chunks in the zone")
	}

	for _, c := range chunks {
		globalIndex, err := extractGlobalIndex(origin, c.Hdr.Name)
		if err != nil {
			return nil, fmt.Errorf("something went wrong during extraction of bloomfilter (global index)")
		}
		bitArray, _, _, err := decodeBloomTxt(c)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		l := uint64(len(*bitArray))
		if l != z.chunkSize {
			return nil, fmt.Errorf("chunksize different to length of this bitarray")
		}
		nr_bits := copy(z.bf.bitArray[globalIndex*l:(globalIndex+1)*l], *bitArray)
		if nr_bits != int(l) {
			return nil, fmt.Errorf("something went wrong during extraction of bloomfilter (copying)")
		}
	}

	log.Infof("\n%s", z.bf.Info())

	return z, nil
}

func decodeBloomTxt(txt_rr *dns.TXT) (*[]bool, uint64, uint64, error) {
	l := len(txt_rr.Txt)
	m, err := strconv.ParseUint((txt_rr.Txt)[l-2], 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	k, err := strconv.ParseUint((txt_rr.Txt)[l-1], 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}

	bloom_string := ""
	for _, s := range txt_rr.Txt[:l-2] {
		bloom_string += s
	}

	decoded_bytes, err := rrutil.FromBase64(bloom_string)
	if err != nil {
		return nil, 0, 0, err
	}

	b := bytesToBits(&decoded_bytes)

	return b, m, k, nil
}
