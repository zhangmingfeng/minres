package router

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	httpserver.RegisterDevDirective("route", "")
	caddy.RegisterPlugin("route", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new route middleware instance.
func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)

	routes, err := routeParse(c)
	if err != nil {
		return err
	}

	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Handler{
			Next:   next,
			Routes: routes,
		}
	})

	return nil
}

func routeParse(c *caddy.Controller) ([]Route, error) {
	var routes []Route
	for c.Next() {
		args := c.RemainingArgs()

		if len(args) != 2 {
			return routes, c.ArgErr()
		}

		route := Route{
			Method:     METHOD_ALL,
			Path:       args[0],
			Controller: args[1],
		}

		for c.NextBlock() {
			switch c.Val() {
			case "method":
				if !c.NextArg() {
					return routes, c.ArgErr()
				}
				route.Method = c.Val()
			case "name":
				if !c.NextArg() {
					return routes, c.ArgErr()
				}
				route.Name = c.Val()
			}
		}

		routes = append(routes, route)
	}
	return routes, nil
}
