package pubsub

import (
	"context"
	"fmt"
)

// Router for pubsub requests
type Router struct {
	routes   map[string]Handler
	notFound Handler
}

// NewRouter create and initialize a router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
		notFound: func(ctx context.Context, req Request) Response {
			return ErrorResponseTo(req, fmt.Errorf("invalid request"))
		},
	}
}

// Add route to the router
func (router *Router) Add(group, entity, method string, handler Handler) {
	router.routes[Route{
		group:  group,
		entity: entity,
		method: method,
	}.String()] = handler
}

// NotFound sets the route handler if no route matches the request
func (router *Router) NotFound(handler Handler) {
	router.notFound = handler
}

// ServeRequest serve a specific request for a given route
func (router *Router) ServeRequest(ctx context.Context, req Request) (resp Response) {
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
	group   string
	entity  string
	method  string
	handler Handler
}

func (r Route) String() string {
	return r.group + "/" + r.entity + "/" + r.method
}

// Handler of requests through pubsub
type Handler func(context.Context, Request) Response
