package webfinger

import "net/http"

// Middleware constant keys
const (
	NoCacheMiddleware     string = "NoCache"
	CorsMiddleware        string = "Cors"
	ContentTypeMiddleware string = "Content-Type"
)

// noCache sets the headers to disable caching
func noCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
}

// jrdSetup sets the content-type
func jrdSetup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/jrd+json")
}
