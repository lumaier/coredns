package bloomfile_nsec5

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/plugin/transfer"
)

func init() { plugin.Register("bloomfile_nsec5", setup) }

func setup(c *caddy.Controller) error {
	zones, err := fileParse(c)
	if err != nil {
		return plugin.Error("bloomfile_nsec5", err)
	}

	f := File{Zones: zones}
	// get the transfer plugin, so we can send notifies and send notifies on startup as well.
	c.OnStartup(func() error {
		t := dnsserver.GetConfig(c).Handler("transfer")
		if t == nil {
			return nil
		}
		f.transfer = t.(*transfer.Transfer) // if found this must be OK.
		go func() {
			for _, n := range zones.Names {
				f.transfer.Notify(n)
			}
		}()
		return nil
	})

	c.OnRestartFailed(func() error {
		t := dnsserver.GetConfig(c).Handler("transfer")
		if t == nil {
			return nil
		}
		go func() {
			for _, n := range zones.Names {
				f.transfer.Notify(n)
			}
		}()
		return nil
	})

	for _, n := range zones.Names {
		z := zones.Z[n]
		c.OnShutdown(z.OnShutdown)
		c.OnStartup(func() error {
			z.StartupOnce.Do(func() { z.Reload(f.transfer) })
			return nil
		})
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})

	return nil
}

func fileParse(c *caddy.Controller) (Zones, error) {
	z := make(map[string]*Zone)
	names := []string{}

	config := dnsserver.GetConfig(c)

	var openErr error
	reload := 1 * time.Minute
	var filepath_vrfkeys string

	for c.Next() {
		// file db.file [zones...]
		if !c.NextArg() {
			return Zones{}, c.ArgErr()
		}
		fileName := c.Val()

		origins := plugin.OriginsFromArgsOrServerBlock(c.RemainingArgs(), c.ServerBlockKeys)
		if !filepath.IsAbs(fileName) && config.Root != "" {
			fileName = filepath.Join(config.Root, fileName)
		}

		reader, err := os.Open(filepath.Clean(fileName))
		if err != nil {
			openErr = err
		}

		for i := range origins {
			z[origins[i]] = NewZone(origins[i], fileName)
			if openErr == nil {
				reader.Seek(0, 0)
				zone, err := Parse(reader, origins[i], fileName, "nil", 0)
				if err != nil {
					return Zones{}, err
				}
				z[origins[i]] = zone
			}
			names = append(names, origins[i])
		}

		for c.NextBlock() {
			switch c.Val() {
			case "reload":
				t := c.RemainingArgs()
				if len(t) < 1 {
					return Zones{}, errors.New("reload duration value is expected")
				}
				d, err := time.ParseDuration(t[0])
				if err != nil {
					return Zones{}, plugin.Error("file", err)
				}
				reload = d
			case "upstream":
				// remove soon
				c.RemainingArgs()
			case "vrf_keys":
				filepath := c.RemainingArgs()
				if len(filepath) != 1 {
					return Zones{}, c.Errf("can only be one argument after %q", "vrf_keys")
				}
				filepath_vrfkeys = filepath[0]
			default:
				return Zones{}, c.Errf("unknown property '%s'", c.Val())
			}
		}

		for i := range origins {
			z[origins[i]].ReloadInterval = reload
			z[origins[i]].Upstream = upstream.New()
			z[origins[i]].vrf_keyfile = filepath_vrfkeys
			err_read := z[origins[i]].ReadKeys(filepath_vrfkeys)
			if err_read != nil {
				return Zones{}, err_read
			}
		}
	}

	if openErr != nil {
		if reload == 0 {
			// reload hasn't been set make this a fatal error
			return Zones{}, plugin.Error("file", openErr)
		}
		log.Warningf("Failed to open %q: trying again in %s", openErr, reload)

	}

	plugin.PrintMemUsage()

	return Zones{Z: z, Names: names}, nil
}
