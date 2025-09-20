package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	t "github.com/easy-comerce/backend/pkg/types"
)

type Route struct {
	method           string
	version          *string
	path             string
	handler          t.HandlerFunc
	routeMiddlewares []t.MiddlewareHandler
	router           *Router
}

type Router struct {
	mux               *http.ServeMux
	globalMiddlewares []t.MiddlewareHandler
	registeredRoutes  map[string]bool
}

func New(mux *http.ServeMux) *Router {
	return &Router{
		mux:               mux,
		globalMiddlewares: make([]t.MiddlewareHandler, 0),
		registeredRoutes:  make(map[string]bool),
	}
}
func (r *Router) GET(path string, handler t.HandlerFunc) *Route {
	return r.addRoute(c.GET, path, handler)
}

func (r *Router) POST(path string, handler t.HandlerFunc) *Route {
	return r.addRoute(c.POST, path, handler)
}

func (r *Router) PUT(path string, handler t.HandlerFunc) *Route {
	return r.addRoute(c.PUT, path, handler)
}

func (r *Router) DELETE(path string, handler t.HandlerFunc) *Route {
	return r.addRoute(c.DELETE, path, handler)
}

func (r *Router) PATCH(path string, handler t.HandlerFunc) *Route {
	return r.addRoute(c.PATCH, path, handler)
}

func (r *Router) addRoute(method, path string, handler t.HandlerFunc) *Route {
	route := &Route{
		method:           method,
		path:             path,
		handler:          handler,
		routeMiddlewares: make([]t.MiddlewareHandler, 0),
		router:           r,
	}

	return route
}

func (route *Route) Version(version string) *Route {
	route.version = &version
	return route
}

func (r *Router) Use(middlewares ...t.MiddlewareHandler) {
	r.globalMiddlewares = append(r.globalMiddlewares, middlewares...)
}

func (route *Route) Use(middlewares ...t.MiddlewareHandler) *Route {
	route.routeMiddlewares = append(route.routeMiddlewares, middlewares...)
	return route
}

func (route *Route) Register() {
	allMiddlewares := make([]t.MiddlewareHandler, 0, len(route.router.globalMiddlewares)+len(route.routeMiddlewares))
	allMiddlewares = append(allMiddlewares, route.router.globalMiddlewares...)
	allMiddlewares = append(allMiddlewares, route.routeMiddlewares...)

	var finalHandler http.Handler = http.HandlerFunc(route.handler)
	for i := 0; i < len(allMiddlewares); i++ {
		finalHandler = allMiddlewares[i](finalHandler)
	}

	route.registerHandler(finalHandler)
}

func (route *Route) registerHandler(handler http.Handler) {
	finalPath := route.buildFinalPathWithVersion()
	routeKey := route.method + " " + finalPath

	if route.router.registeredRoutes[routeKey] {
		if config.GetActiveProfile() != c.EnvProd {
			panic(fmt.Sprintf("Duplicate route detected: %s", routeKey))
		}
		logger.Logger.Warn("WARNING: Duplicate route detected. Overwriting existing route", "routeKey", routeKey)
	}
	route.router.registeredRoutes[routeKey] = true
	route.router.mux.Handle(routeKey, handler)
}

func (route *Route) buildFinalPathWithVersion() string {
	path := route.path

	if route.version != nil {
		versionPrefix := c.VersionPrefix
		pathStartsWithV := strings.HasPrefix(path, versionPrefix)
		pathIsLongEnoughToHaveVersionNumber := len(path) > len(versionPrefix)
		characterAfterV := path[len(versionPrefix)]
		characterAfterVIsDigit := characterAfterV >= '0' && characterAfterV <= '9'

		pathHasVersionPrefix := pathStartsWithV && pathIsLongEnoughToHaveVersionNumber && characterAfterVIsDigit
		if pathHasVersionPrefix {
			pathSegments := strings.SplitN(path, "/", 3)
			pathHasRouteAfterVersion := len(pathSegments) >= 3

			if pathHasRouteAfterVersion {
				routeAfterVersion := pathSegments[2]
				path = "/" + routeAfterVersion
			}
		}

		finalVersionedPath := "/" + *route.version + path
		path = finalVersionedPath
	}

	return path
}
