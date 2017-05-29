package pubsub

import (
	"context"
	"fmt"
)

// Router for pubsub requests
type Router struct {
	routes   map[string]Endpoint
	notFound Endpoint
}

// NewRouter create and initialize a router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Endpoint),
		notFound: func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			err = fmt.Errorf("invalid request")
			return
		},
	}
}

// Add route to the router
func (router *Router) Add(group, entity, method string, ep Endpoint) {
	router.routes[Route{
		group:  group,
		entity: entity,
		method: method,
	}.String()] = ep
}

// NotFound sets the route handler if no route matches the request
func (router *Router) NotFound(ep Endpoint) {
	router.notFound = ep
}

// ServeRequest serve a specific request for a given route
func (router *Router) ServeRequest(ctx context.Context, req Request) (interface{}, error) {
	routeToFind := Route{
		group:  req.Group,
		entity: req.Entity,
		method: req.Method,
	}
	if handler, ok := router.routes[routeToFind.String()]; ok {
		return handler(ctx, req)
	}
	return router.notFound(ctx, req)
}

// Route defines a rule for routing
type Route struct {
	group  string
	entity string
	method string
}

func (r Route) String() string {
	return r.group + "/" + r.entity + "/" + r.method
}

// Endpoint of requests through pubsub
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
