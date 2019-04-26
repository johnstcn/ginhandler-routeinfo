package ginrouteinfo

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// RoutePathKey is the key under which the matched route path is stored.
const RoutePathKey = "route_path"

// A RoutesInfoProvider is a function that returns gin.RoutesInfo.
type RoutesInfoProvider func() gin.RoutesInfo

// WithRoutePath returns a new gin.HandlerFunc that stores the path of the request in the context under the context key `route_path`.
// This replaces raw parameter values with their wildcard names. For example, given a router path `/foo/:bar`, for a request
// `GET /foo/baz` the stored value of `route_path` will be `/foo/:bar`.
// This is useful for logging all request URL paths without overly inflating metric cardinality.
// If the HandlerName of the incoming request is not known, `route_path` will be set to an empty string.
// Most of the implementation is inspired by https://github.com/gin-gonic/gin/issues/748
// While there are some open PRs to add this to Gin, none of them are merged at time of writing.
func WithRoutePath(routes RoutesInfoProvider) gin.HandlerFunc {
	rl := newRouteLookuper(routes)
	return func(g *gin.Context) {
		var path string
		var found bool

		// The first time rl calls its RoutesInfoProvider, all routes may not be present.
		// If we can't find the handler name already there, update once and return that result.
		if path, found = rl.Lookup(g.HandlerName()); !found {
			rl.Update()
			path, found = rl.Lookup(g.HandlerName())
		}

		g.Set(RoutePathKey, path)
		g.Next()
	}
}

func newRouteLookuper(routes RoutesInfoProvider) *routeLookuper {
	rl := routeLookuper{
		routes: routes,
	}
	rl.lock.cached = make(map[string]string)
	return &rl
}

type routeLookuper struct {
	routes RoutesInfoProvider
	lock   struct {
		sync.RWMutex
		cached map[string]string
	}
}

func (r *routeLookuper) Lookup(handler string) (path string, found bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	path, found = r.lock.cached[handler]
	return path, found
}

func (r *routeLookuper) Update() {
	cache := make(map[string]string)
	for _, ri := range r.routes() {
		cache[ri.Handler] = ri.Path
	}
	r.lock.Lock()
	r.lock.cached = cache
	r.lock.Unlock()
}
