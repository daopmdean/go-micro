package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRouteExist(t *testing.T) {
	tApp := Config{}

	routes := tApp.routes()
	chiRoutes := routes.(chi.Router)

	needRoutes := []string{"/authenticate"}

	for _, need := range needRoutes {
		routeExist(t, chiRoutes, need)
	}
}

func routeExist(t *testing.T, routes chi.Router, route string) {
	found := false

	chi.Walk(routes, func(
		method, foundRoute string,
		handler http.Handler,
		middlewares ...func(http.Handler) http.Handler,
	) error {
		if route == foundRoute {
			found = true
		}

		return nil
	})

	if !found {
		t.Errorf("route not found: %s", route)
	}
}
