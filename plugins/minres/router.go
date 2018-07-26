package minres

import (
	"github.com/gorilla/mux"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net/http"
)

type Handler struct {
	Next   httpserver.Handler
	Router *mux.Router
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
