package router

import (
	"net/http"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/gorilla/mux"
)

type Handler struct {
	Next   httpserver.Handler
	Router *mux.Router
	HOST   string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var match mux.RouteMatch
	var handler http.Handler
	if h.Router.Match(r, &match) {
		handler = match.Handler
	}

	if handler != nil || (handler != nil && match.MatchErr == mux.ErrMethodMismatch) {
		h.Router.ServeHTTP(w, r)
		return 200, nil
	}

	return h.Next.ServeHTTP(w, r)
}
