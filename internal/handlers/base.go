package handlers_

import (
	"net/http"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var routes []Route

func NewRoute(method string, path string, handler http.HandlerFunc) Route {
	return Route{method, path, handler}
}
func RegisterRoute(route Route) {
	routes = append(routes, route)
}
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, route := range routes {
		if route.Method == r.Method && route.Path == r.URL.Path {
			route.Handler(w, r)
			found = true
			break
		}
	}
	if !found {
		http.NotFound(w, r)
	}
}
func InitBaseRoutes() {
	routes = append(routes, NewRoute("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
}
