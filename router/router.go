package router

import (
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	response "xiview.com/http/response/json"
	"github.cn/mholt/caddy/http/controller"
	"regexp"
	"errors"
	"fmt"
)

type Handler struct {
	Next   httpserver.Handler
	Routes []Route
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	for _, route := range h.Routes {
		hpath := httpserver.Path(r.URL.Path)
		params, err := MathchRoute(string(hpath), route.Path)
		if err != nil {
			if strings.HasPrefix(string(hpath), "/") {
				// this is a normal-looking path, and it doesn't match; try next rule
				continue
			}
			hpath = httpserver.Path("/" + string(hpath)) // prepend leading slash
			if !strings.EqualFold(string(hpath), route.Path) {
				// even after fixing the request path, it still doesn't match; try next rule
				continue
			}
		}
		allowedMethod := route.Method
		if !strings.Contains(allowedMethod, METHOD_ALL) && !strings.Contains(allowedMethod, r.Method) {
			jsonRes := response.Json{
				Code: http.StatusForbidden,
				Msg:  "not allowed!",
			}
			return jsonRes.ResponseBody(w)
		}
		controller, ok := registerController[route.Controller]
		if ok {
			return controller.Handle(w, r, params)
		} else {
			jsonRes := response.Json{
				Code: http.StatusForbidden,
				Msg:  fmt.Sprintf("action [%s] not found!", route.Controller),
			}
			return jsonRes.ResponseBody(w)
		}
	}
	return h.Next.ServeHTTP(w, r)
}

func RegisterController(name string, controller controller.Controller) {
	if name == "" {
		panic("controller must has a name")
	}
	if _, ok := registerController[name]; !ok {
		registerController[name] = controller
	}
}

func MathchRoute(reqPath string, routePath string) (map[string]interface{}, error) {
	var params = make(map[string]interface{})
	if strings.EqualFold(reqPath, routePath) {
		return nil, nil
	}
	match, _ := regexp.MatchString("{.*}", routePath)
	if !match {
		return nil, errors.New("not matched")
	}
	reqPath = strings.TrimRight(strings.TrimLeft(reqPath, "/"), "/")
	routePath = strings.TrimRight(strings.TrimLeft(routePath, "/"), "/")
	reqPaths := strings.Split(reqPath, "/")
	routePaths := strings.Split(routePath, "/")
	if len(reqPaths) != len(routePaths) {
		return nil, errors.New("not matched")
	}
	for i, item := range routePaths {
		paramed, _ := regexp.MatchString("{.*}", item)
		if !paramed {
			if item != reqPaths[i] {
				return nil, errors.New("not matched")
			}
		} else {
			item = strings.Trim(strings.Trim(item, "{"), "}")
			params[item] = reqPaths[i]
		}
	}
	return params, nil
}

type Route struct {
	Method     string
	Name       string
	Path       string
	Controller string
}

const (
	METHOD_ALL  = "ALL"
	METHOD_GET  = "GET"
	METHOD_POST = "POST"
)

var (
	registerController = make(map[string]controller.Controller)
)
