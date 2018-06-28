package router

import (
	"github.com/mholt/caddy"
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/gorilla/mux"
	"strings"
)

func init() {
	caddy.RegisterPlugin("route", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)
	router, err := routeParse(c)
	if err != nil {
		return err
	}

	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Handler{
			Next:   next,
			Router: router,
		}
	})

	return nil
}

func routeParse(c *caddy.Controller) (*mux.Router, error) {
	r := mux.NewRouter()
	for c.Next() {
		args := c.RemainingArgs()

		if len(args) != 1 {
			return r, c.ArgErr()
		}

		path := args[0]
		var method = ""
		var action string

		for c.NextBlock() {
			switch c.Val() {
			case "action":
				if !c.NextArg() {
					return r, c.ArgErr()
				}
				action = c.Val()
			case "method":
				if !c.NextArg() {
					return r, c.ArgErr()
				}
				method = c.Val()
			}
		}

		if handle, ok := registerController[action]; ok {
			if len(method) > 0 {
				methods := strings.Split(strings.Replace(method, " ", "", -1), ",")
				r.HandleFunc(path, handle).Name(path).Methods(methods...)
			} else {
				r.HandleFunc(path, handle).Name(path)
			}

		}
	}
	return r, nil
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
)
