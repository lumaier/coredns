package bloomsec_nsec5

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("bloomsec_nsec5", setup) }

func setup(c *caddy.Controller) error {
	sign, err := parse(c)
	if err != nil {
		return plugin.Error("bloomfile_nsec5", err)
	}

	c.OnStartup(sign.OnStartup)
	c.OnStartup(func() error {
		for _, signer := range sign.signers {
			go signer.refresh(durationRefreshHours)
		}
		return nil
	})
	c.OnShutdown(func() error {
		for _, signer := range sign.signers {
			close(signer.stop)
		}
		return nil
	})

	// Don't call AddPlugin, *bloomsec* is not a plugin.
	return nil
}

func parse(c *caddy.Controller) (*Sign, error) {
	sign := &Sign{}
	config := dnsserver.GetConfig(c)

	for c.Next() {
		if !c.NextArg() {
			return nil, c.ArgErr()
		}
		dbfile := c.Val()
		if !filepath.IsAbs(dbfile) && config.Root != "" {
			dbfile = filepath.Join(config.Root, dbfile)
		}

		origins := plugin.OriginsFromArgsOrServerBlock(c.RemainingArgs(), c.ServerBlockKeys)
		signers := make([]*Signer, len(origins))
		for i := range origins {
			signers[i] = &Signer{
				dbfile:      dbfile,
				origin:      origins[i],
				jitterIncep: time.Duration(float32(durationInceptionJitter) * rand.Float32()),
				jitterExpir: time.Duration(float32(durationExpirationDayJitter) * rand.Float32()),
				directory:   "/var/lib/coredns",
				stop:        make(chan struct{}),
				signedfile:  fmt.Sprintf("db.%ssigned", origins[i]), // origins[i] is a fqdn, so it ends with a dot, hence %ssigned.
				fp_prob:     0.001,
				chunkSize:   uint64(24 * 1200),
			}
		}

		for c.NextBlock() {
			switch c.Val() {
			case "key":
				pairs, err := keyParse(c)
				if err != nil {
					return sign, err
				}
				for i := range signers {
					for _, p := range pairs {
						p.Public.Header().Name = signers[i].origin
					}
					signers[i].keys = append(signers[i].keys, pairs...)
				}
			case "directory":
				dir := c.RemainingArgs()
				if len(dir) == 0 || len(dir) > 1 {
					return sign, fmt.Errorf("can only be one argument after %q", "directory")
				}
				if !filepath.IsAbs(dir[0]) && config.Root != "" {
					dir[0] = filepath.Join(config.Root, dir[0])
				}
				for i := range signers {
					signers[i].directory = dir[0]
					signers[i].signedfile = fmt.Sprintf("db.%ssigned", signers[i].origin)
				}
			case "fp_rate":
				fprate := c.RemainingArgs()
				if len(fprate) != 1 {
					return sign, fmt.Errorf("can only be one argument after %q", "FPrate")
				}
				f, err := strconv.ParseFloat(fprate[0], 64)
				if err != nil {
					return sign, err
				}
				for i := range signers {
					signers[i].fp_prob = f
				}
			case "chunksize":
				chunksize := c.RemainingArgs()
				if len(chunksize) != 1 {
					return sign, fmt.Errorf("can only be one argument after %q", "chunksize")
				}
				c, err := strconv.ParseUint(chunksize[0], 10, 64)
				if err != nil {
					return sign, err
				}
				if c%bitsEncoded != 0 {
					return sign, fmt.Errorf("chunksize needs to be a multiple of %d", bitsEncoded)
				}
				for i := range signers {
					signers[i].chunkSize = c
				}
			default:
				return nil, c.Errf("unknown property '%s'", c.Val())
			}
		}
		sign.signers = append(sign.signers, signers...)
	}

	return sign, nil
}
