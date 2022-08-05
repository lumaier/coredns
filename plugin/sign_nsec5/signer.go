package sign_nsec5

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/ProtonMail/go-ecvrf/ecvrf"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/file_nsec5"
	"github.com/coredns/coredns/plugin/file_nsec5/tree"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("sign_nsec5")

// Signer holds the data needed to sign a zone file.
type Signer struct {
	keys        []Pair
	origin      string
	dbfile      string
	directory   string
	jitterIncep time.Duration
	jitterExpir time.Duration

	signedfile string
	stop       chan struct{}

	vrf_privkey *ecvrf.PrivateKey
	vrf_pubkey  *ecvrf.PublicKey
}

// Sign signs a zone file according to the parameters in s.
// Additionally it creates the Bloomfilter and the corresponding TXT records.
func (s *Signer) Sign(now time.Time) (*file_nsec5.Zone, error) {
	rd, err := os.Open(s.dbfile)
	if err != nil {
		return nil, err
	}

	z, err := Parse(rd, s.origin, s.dbfile)
	if err != nil {
		return nil, err
	}

	// mttl := z.Apex.SOA.Minttl
	ttl := z.Apex.SOA.Header().Ttl
	inception, expiration := lifetime(now, s.jitterIncep, s.jitterExpir)
	z.Apex.SOA.Serial = uint32(now.Unix())

	for _, pair := range s.keys {
		pair.Public.Header().Ttl = ttl // set TTL on key so it matches the RRSIG.
		z.Insert(pair.Public)
		z.Insert(pair.Public.ToDS(dns.SHA1).ToCDS())
		z.Insert(pair.Public.ToDS(dns.SHA256).ToCDS())
		z.Insert(pair.Public.ToCDNSKEY())
	}

	nsec5_names := []string{}

	for _, pair := range s.keys {
		rrsig, err := pair.signRRs([]dns.RR{z.Apex.SOA}, s.origin, ttl, inception, expiration)
		if err != nil {
			return nil, err
		}
		z.Insert(rrsig)
		// NS apex may not be set if RR's have been discarded because the origin doesn't match.
		if len(z.Apex.NS) > 0 {
			rrsig, err = pair.signRRs(z.Apex.NS, s.origin, ttl, inception, expiration)
			if err != nil {
				return nil, err
			}
			z.Insert(rrsig)
		}
	}

	err = z.AuthWalk(func(e *tree.Elem, zrrs map[uint16][]dns.RR, auth bool) error {
		if !auth {
			return nil
		}
		nsec5_names = append(nsec5_names, e.Name())
		return nil
	})

	//////////////////////////// NSEC5 ///////////////////////////////////////////////////

	// create NSEC5 and NSEC5PROOF
	nsec5_values := make([]VRF_output, len(nsec5_names)) // holds: vrf hash, vrf proof, domain name

	wg := sync.WaitGroup{}

	wg.Add(len(nsec5_names))

	for i, name := range nsec5_names {
		go func(i int, x string) {
			defer wg.Done()
			pi, hash, _ := s.vrf_privkey.Prove([]byte(x))
			if err != nil {
				panic(err)
			}
			nsec5_values[i] = VRF_output{
				hash:   toBase64(hash),
				proof:  toBase64(pi),
				domain: x,
			}
		}(i, name)
	}

	wg.Wait()

	// sort the vrf hash outputs
	sort.Slice(nsec5_values, func(i, j int) bool {
		return nsec5_values[i].hash < nsec5_values[j].hash
	})

	length_nsec5 := len(nsec5_values)
	proofs := []*dns.TXT{}

	// for every two consecutive domain names insert an NSEC5 and NSEC5PROOF
	for i, x := range nsec5_values {
		r1, r2 := NSEC5(x.domain, x.hash, x.proof, nsec5_values[(i+1)%length_nsec5].hash, ttl)
		z.Insert(r1)
		proofs = append(proofs, r2) // FIXME: the proofs are inserted later as they do not have to be signed?
	}

	//////////////////////////////////////////////////////////////////////////////////////

	// sign the records we are authoritative for
	err = z.AuthWalk(func(e *tree.Elem, zrrs map[uint16][]dns.RR, auth bool) error {
		if !auth {
			return nil
		}

		for t, rrs := range zrrs {
			// RRSIGs are not signed and NS records are not signed because we are never authoritative for them.
			// The zone's apex nameservers records are not kept in this tree and are signed separately.
			if t == dns.TypeRRSIG || t == dns.TypeNS {
				continue
			}
			for _, pair := range s.keys {
				rrsig, err := pair.signRRs(rrs, s.origin, rrs[0].Header().Ttl, inception, expiration)
				if err != nil {
					return err
				}
				e.Insert(rrsig)
			}
		}
		return nil
	})

	for _, x := range proofs {
		z.Insert(x)
	}

	return z, err
}

// resign checks if the signed zone exists, or needs resigning.
func (s *Signer) resign() error {
	signedfile := filepath.Join(s.directory, s.signedfile)
	rd, err := os.Open(filepath.Clean(signedfile))
	if err != nil && os.IsNotExist(err) {
		return err
	}

	now := time.Now().UTC()
	return resign(rd, now)
}

// resign will scan rd and check the signature on the SOA record. We will resign on the basis
// of 2 conditions:
// * either the inception is more than 6 days ago, or
// * we only have 1 week left on the signature
//
// All SOA signatures will be checked. If the SOA isn't found in the first 100
// records, we will resign the zone.
func resign(rd io.Reader, now time.Time) (why error) {
	zp := dns.NewZoneParser(rd, ".", "resign")
	zp.SetIncludeAllowed(true)
	i := 0

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		if err := zp.Err(); err != nil {
			return err
		}

		switch x := rr.(type) {
		case *dns.RRSIG:
			if x.TypeCovered != dns.TypeSOA {
				continue
			}
			incep, _ := time.Parse("20060102150405", dns.TimeToString(x.Inception))
			// If too long ago, resign.
			if now.Sub(incep) >= 0 && now.Sub(incep) > durationResignDays {
				return fmt.Errorf("inception %q was more than: %s ago from %s: %s", incep.Format(timeFmt), durationResignDays, now.Format(timeFmt), now.Sub(incep))
			}
			// Inception hasn't even start yet.
			if now.Sub(incep) < 0 {
				return fmt.Errorf("inception %q date is in the future: %s", incep.Format(timeFmt), now.Sub(incep))
			}

			expire, _ := time.Parse("20060102150405", dns.TimeToString(x.Expiration))
			if expire.Sub(now) < durationExpireDays {
				return fmt.Errorf("expiration %q is less than: %s away from %s: %s", expire.Format(timeFmt), durationExpireDays, now.Format(timeFmt), expire.Sub(now))
			}
		}
		i++
		if i > 100 {
			// 100 is a random number. A SOA record should be the first in the zonefile, but RFC 1035 doesn't actually mandate this. So it could
			// be 3rd or even later. The number 100 looks crazy high enough that it will catch all weird zones, but not high enough to keep the CPU
			// busy with parsing all the time.
			return fmt.Errorf("no SOA RRSIG found in first 100 records")
		}
	}

	return nil
}

func signAndLog(s *Signer, why error) {
	now := time.Now().UTC()
	z, err := s.Sign(now)
	log.Infof("Signing %q because %s", s.origin, why)
	if err != nil {
		log.Warningf("Error signing %q with key tags %q in %s: %s, next: %s", s.origin, keyTag(s.keys), time.Since(now), err, now.Add(durationRefreshHours).Format(timeFmt))
		return
	}

	if err := s.write(z); err != nil {
		log.Warningf("Error signing %q: failed to move zone file into place: %s", s.origin, err)
		return
	}
	log.Infof("Successfully signed zone %q in %q with key tags %q and %d SOA serial, elapsed %f, next: %s", s.origin, filepath.Join(s.directory, s.signedfile), keyTag(s.keys), z.Apex.SOA.Serial, time.Since(now).Seconds(), now.Add(durationRefreshHours).Format(timeFmt))

	plugin.PrintMemUsage()
}

// refresh checks every val if some zones need to be resigned.
func (s *Signer) refresh(val time.Duration) {
	tick := time.NewTicker(val)
	defer tick.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-tick.C:
			why := s.resign()
			if why == nil {
				continue
			}
			signAndLog(s, why)
		}
	}
}

func lifetime(now time.Time, jitterInception, jitterExpiration time.Duration) (uint32, uint32) {
	incep := uint32(now.Add(durationSignatureInceptionHours).Add(jitterInception).Unix())
	expir := uint32(now.Add(durationSignatureExpireDays).Add(jitterExpiration).Unix())
	return incep, expir
}
