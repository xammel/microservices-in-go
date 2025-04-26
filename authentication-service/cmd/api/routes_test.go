package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config {
		
	}

	testRoutes := testApp.routes()

	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}

}

func routeExists(t *testing.T, router chi.Router, route string) {
	found := false

	// Check if the route exists in the router
	_ = chi.Walk(router, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route {
			found = true
			return nil
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}