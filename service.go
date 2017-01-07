package webfinger

import (
	"net/http"

	"github.com/captncraig/cors"
)

// Service is the webfinger service containing the required
// HTTP handlers and defaults for webfinger implementations.
type Service struct {

	// PreHandlers are invoked at the start of each HTTP method, used to
	// setup things like CORS, caching, etc. You can delete or replace
	// a handler by setting service.PreHandlers[name] = nil
	PreHandlers map[string]http.Handler

	// NotFoundHandler is the handler to invoke when a URL is not matched. It does NOT
	// handle the case of a non-existing users unless your Resolver.DummyUser returns
	// an error that matches Resolver.IsNotFoundError(err) == true.
	NotFoundHandler http.Handler

	// MethodNotSupportedHandler is the handler invoked when an unsupported
	// method is called on the webfinger HTTP service.
	MethodNotSupportedHandler http.Handler

	// MalformedRequestHandler is the handler invoked if the request routes
	// but is malformed in some way. The default behavior is to return 400 BadRequest,
	// per the webfinger specification
	MalformedRequestHandler http.Handler

	// NoTLSHandler is the handler invoked if the request is not
	// a TLS request. The default behavior is to redirect to the TLS
	// version of the URL, per the webfinger specification. Setting
	// this to nil will allow nonTLS connections, but that is not advised.
	NoTLSHandler http.Handler

	// ErrorHandler is the handler invoked when an error is called. The request
	// context contains the error in the webfinger.ErrorKey and can be fetched
	// via webfinger.ErrorFromContext(ctx)
	ErrorHandler http.Handler

	// Resolver is the interface for resolving user details
	Resolver Resolver
}

// Default creates a new service with the default registered handlers
func Default(ur Resolver) *Service {
	var c = cors.Default()

	s := &Service{}
	s.Resolver = ur
	s.ErrorHandler = http.HandlerFunc(s.defaultErrorHandler)
	s.NotFoundHandler = http.HandlerFunc(s.defaultNotFoundHandler)
	s.MethodNotSupportedHandler = http.HandlerFunc(s.defaultMethodNotSupportedHandler)
	s.MalformedRequestHandler = http.HandlerFunc(s.defaultMalformedRequestHandler)
	s.NoTLSHandler = http.HandlerFunc(s.defaultNoTLSHandler)

	s.PreHandlers = make(map[string]http.Handler)
	s.PreHandlers[NoCacheMiddleware] = http.HandlerFunc(noCache)
	s.PreHandlers[CorsMiddleware] = http.HandlerFunc(c.HandleRequest)
	s.PreHandlers[ContentTypeMiddleware] = http.HandlerFunc(jrdSetup)

	return s
}
