package router

import (
	"fmt"
	"net/http"
	"strings"

	c "github.com/easy-comerce/backend/pkg/constants"
	m "github.com/easy-comerce/backend/pkg/middleware"
)

type Route struct {
	method           string
	version          *string
	path             string
	handler          m.HandlerFunc
	routeMiddlewares []m.MiddlewareHandler
	router           *Router
}

type Router struct {
	mux               *http.ServeMux
	globalMiddlewares []m.MiddlewareHandler
	routeMiddlewares  []m.MiddlewareHandler
	handlers          map[string]http.Handler
}

func New(mux *http.ServeMux) *Router {
	return &Router{
		mux:               mux,
		globalMiddlewares: make([]m.MiddlewareHandler, 0),
		routeMiddlewares:  make([]m.MiddlewareHandler, 0),
		handlers:          make(map[string]http.Handler),
	}
}

func (r *Router) GET(path string, handler m.HandlerFunc) *Route {
	return r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler m.HandlerFunc) *Route {
	return r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler m.HandlerFunc) *Route {
	return r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler m.HandlerFunc) *Route {
	return r.addRoute(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler m.HandlerFunc) *Route {
	return r.addRoute(http.MethodPatch, path, handler)
}

func (r *Router) addRoute(method, path string, handler m.HandlerFunc) *Route {
	route := &Route{
		method:           method,
		path:             path,
		handler:          handler,
		routeMiddlewares: make([]m.MiddlewareHandler, 0),
		router:           r,
	}

	return route
}

func (route *Route) Version(version string) *Route {
	route.version = &version
	return route
}

func (r *Router) Use(middlewares ...m.MiddlewareHandler) http.Handler {
	r.globalMiddlewares = append(r.globalMiddlewares, middlewares...)

	var handler http.Handler = r.mux
	for i := len(r.globalMiddlewares) - 1; i >= 0; i-- {
		handler = r.globalMiddlewares[i](handler)
	}
	return handler
}

func (r *Router) GetGlobalMiddlewares() []m.MiddlewareHandler {
	return r.globalMiddlewares
}

func (route *Route) Use(middlewares ...m.MiddlewareHandler) *Route {
	route.routeMiddlewares = append(route.routeMiddlewares, middlewares...)
	return route
}

func (route *Route) Register() {
	allMiddlewares := make([]m.MiddlewareHandler, 0, len(route.router.globalMiddlewares)+len(route.routeMiddlewares))
	allMiddlewares = append(allMiddlewares, route.router.globalMiddlewares...)
	allMiddlewares = append(allMiddlewares, route.routeMiddlewares...)

	var finalHandler http.Handler = http.HandlerFunc(route.handler)
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		finalHandler = allMiddlewares[i](finalHandler)
	}
	route.registerHandler(finalHandler)
}

func (route *Route) registerHandler(handler http.Handler) {
	finalPath := route.buildFinalPathWithVersion()
	routeKey := route.method + " " + finalPath

	if _, exists := route.router.handlers[routeKey]; exists {
		panic(fmt.Sprintf("❌ DUPLICATE ROUTE DETECTED: %s", routeKey))
	}

	route.router.handlers[routeKey] = handler
	dispatcherNotYetRegistered := true

	for key := range route.router.handlers {
		if strings.HasPrefix(key, route.method+" "+finalPath) || strings.HasSuffix(key, " "+finalPath) {
			parts := strings.Split(key, " ")

			if len(parts) == 2 && parts[1] == finalPath && parts[0] != route.method {
				dispatcherNotYetRegistered = false
				break
			}
		}
	}

	if dispatcherNotYetRegistered {
		dispatcher := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerKey := r.Method + " " + finalPath
			methodHandler, exists := route.router.handlers[handlerKey]

			if !exists {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			methodHandler.ServeHTTP(w, r)
		})

		route.router.mux.Handle(finalPath, dispatcher)
	}
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
