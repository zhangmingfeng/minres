package minres

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/zhangmingfeng/minres/plugins/minres/weed"
	"net/http"
)

func init() {
	httpserver.RegisterDevDirective("minres", "")
	caddy.RegisterPlugin("minres", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

var host string
var WeedClient *weed.Client

type Config struct {
	WdMaster  string
	WdFilers  []string
	CachePath string
}

func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)
	router, minresCfg, err := parse(c)
	if err != nil {
		return err
	}
	for k, v := range registerController {
		router.HandleFunc(k, v).Name(k)
	}
	routerHandle = Handler{
		Router: router,
	}
	host = cfg.Addr.String()
	WeedClient = weed.NewClient(minresCfg.WdMaster, minresCfg.CachePath, minresCfg.WdFilers...)
	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		routerHandle.Next = next
		return routerHandle
	})

	return nil
}

func Url(name string, params ...string) string {
	route := routerHandle.Router.Get(name)
	if route == nil {
		return ""
	}
	url, _ := route.URL(params...)
	return fmt.Sprintf("%s%s", host, url.String())
}

func parse(c *caddy.Controller) (*mux.Router, *Config, error) {
	r := mux.NewRouter()
	cfg := &Config{WdFilers: []string{}}
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "wd_master":
				if !c.NextArg() {
					return r, cfg, c.ArgErr()
				}
				cfg.WdMaster = c.Val()
			case "wd_filer":
				if !c.NextArg() {
					return r, cfg, c.ArgErr()
				}
				cfg.WdFilers = append(cfg.WdFilers, c.Val())
			case "cache_path":
				if !c.NextArg() {
					return r, cfg, c.ArgErr()
				}
				cfg.CachePath = c.Val()
			}
		}
	}
	return r, cfg, nil
}

func RegisterController(name string, handler http.HandlerFunc) {
	if len(name) <= 0 {
		panic("controller must has a name")
	}
	if _, ok := registerController[name]; !ok {
		registerController[name] = handler
	}
}

var (
	registerController = make(map[string]http.HandlerFunc)
	routerHandle       Handler
)
