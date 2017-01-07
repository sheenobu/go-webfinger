[![GoDoc](https://godoc.org/github.com/sheenobu/go-webfinger?status.svg)](https://godoc.org/github.com/sheenobu/go-webfinger)

# go-webfinger

go-webfinger is a golang webfinger server implementation.

## Usage

`webfinger.Service` is implemented as a net/http handler, which means
usage is simply registering the object with your http service.

Using the webfinger service as the main ServeHTTP:

```go
myResolver = ...
wf := webfinger.Default(myResolver{})
wf.NotFoundHandler = // the rest of your app
http.ListenAndService(":8080", wf)
```

Using the webfinger service as a route on an existing HTTP router:

```go
myResolver = ...
wf := webfinger.Default(myResolver{})
http.Handle(webfinger.WebFingerPath, http.HandlerFunc(wf.Webfinger))
http.ListenAndService(":8080", nil)
```

## Defaults

The webfinger service is installed with a few defaults. Some of these
defaults ensure we stick closely to the webfinger specification (tls-only, CORS, Content-Type)
and other defaults are simply useful for development (no-cache)

The full list of defaults can be found in the godoc for `webfinger.Service`. They are exposed
as public variables which can be overriden.

`PreHandlers` are the list of preflight HTTP handlers to run. You can add your own via `wf.PreHandlers["my-custom-name"] = ...`, however,
execution order is not guaranteed.

### TLS-Only

Handler which routes to the TLS version of the page. Disable via `wf.NoTLSHandler = nil`.

### No-Cache

A PreFlight handler which sets no-cache headers on anything under `/.well-known/webfinger`. Disable or override via `wf.PreHandlers[webfinger.NoCacheMiddleware]`

### Content Type as application/jrd+json

A PreFlight handler which sets the Content-Type to `application/jrd+json`. Disable or override via `wf.PreHandlers[webfinger.ContentTypeMiddleware]`.

### CORS 

A PreFlight handler which adds the CORS headers. Disable or override via `wf.PreHandlers[webfinger.CorsMiddleware].`
