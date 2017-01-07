package webfinger

// The Resolver is how the webfinger service looks up a user or resource. The resolver
// must be provided by the developer using this package, as each webfinger
// service may be exposing a different set or users or resources or services.
type Resolver interface {

	// FindUser finds the user given the username and hostname.
	FindUser(username string, hostname string, r []Rel) (*Resource, error)

	// DummyUser allows us to return a dummy user to avoid user-enumeration via webfinger 404s. This
	// can be done in the webfinger code itself but then it would be obvious which users are real
	// and which are not real via differences in how the implementation works vs how
	// the general webfinger code works. This does not match the webfinger specification
	// but is an extra precaution. Returning a NotFound error here will
	// keep the webfinger 404 behavior.
	DummyUser(username string, hostname string, r []Rel) (*Resource, error)

	// IsNotFoundError returns true if the given error is a not found error.
	IsNotFoundError(err error) bool
}
