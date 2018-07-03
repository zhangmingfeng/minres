package seaweedfs

import (
	"github.com/mholt/caddy"
	"errors"
)

func init() {
	caddy.RegisterPlugin("weedfs", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	weedfsCfg, err := weedfsConfigParse(c)
	if err != nil {
		return err
	}
	defaultClient = NewClient(weedfsCfg.Master, weedfsCfg.Filers...)
	return nil
}

func weedfsConfigParse(c *caddy.Controller) (*WeedfsConfig, error) {
	weedfsCfg := &WeedfsConfig{Filers: []string{}}
	for c.Next() {
		if len(weedfsCfg.Master) > 0 {
			return weedfsCfg, errors.New("duplication weedfs config")
		}
		for c.NextBlock() {
			switch c.Val() {
			case "master":
				if !c.NextArg() {
					return weedfsCfg, c.ArgErr()
				}
				weedfsCfg.Master = c.Val()
			case "filer":
				if !c.NextArg() {
					return weedfsCfg, c.ArgErr()
				}
				weedfsCfg.Filers = append(weedfsCfg.Filers, c.Val())
			}
		}
	}
	return weedfsCfg, nil
}

type WeedfsConfig struct {
	Master string
	Filers []string
}
