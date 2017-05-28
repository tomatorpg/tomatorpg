package pubsub

import "context"

// Router for pubsub requests
type Router struct {
	routes   []Route
	notFound Handler
}

// Add route to the router
func (router *Router) Add(group, entity, method string, handler Handler) {
	router.routes = append(
		router.routes,
		Route{
			group:   group,
			entity:  entity,
			method:  method,
			handler: handler,
		},
	)
}

// NotFound sets the route handler if no route matches the request
func (router *Router) NotFound(handler Handler) {
	router.notFound = handler
}

// Route defines a rule for routing
type Route struct {
	group   string
	entity  string
	method  string
	handler Handler
}

// Handler of requests through pubsub
type Handler func(context.Context, Request) Response
