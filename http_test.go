package webfinger

import (
	"bytes"
	"net/http"
	"net/url"
	"sort"
	"testing"

	"reflect"

	"crypto/tls"

	"strings"

	"github.com/pkg/errors"
)

type dummyUserResolver struct {
}

func (d *dummyUserResolver) FindUser(username string, hostname string, rel []Rel) (*Resource, error) {
	if username == "hello" {
		if len(rel) == 2 && rel[0] == "x" && rel[1] == "y" {
			return &Resource{
				Links: []Link{
					Link{
						HRef: string(rel[0]),
						Rel:  string(rel[0]),
					},
					Link{
						HRef: string(rel[1]),
						Rel:  string(rel[1]),
					},
				},
			}, nil
		}
		return &Resource{
			Links: []Link{
				Link{
					HRef: string("x"),
					Rel:  string("x"),
				},
				Link{
					HRef: string("y"),
					Rel:  string("y"),
				},
				Link{
					HRef: string("z"),
					Rel:  string("z"),
				},
			},
		}, nil
	}

	return nil, errors.New("User not found")
}

func (d *dummyUserResolver) DummyUser(username string, hostname string, rel []Rel) (*Resource, error) {
	return nil, errors.New("User not found")
}

func (d *dummyUserResolver) IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == "User not found" {
		return true
	}
	if errors.Cause(err).Error() == "User not found" {
		return true
	}
	return false
}

type dummyResponseWriter struct {
	bytes.Buffer

	code    int
	headers http.Header
}

func (d *dummyResponseWriter) Header() http.Header {
	if d.headers == nil {
		d.headers = make(http.Header)
	}
	return d.headers
}

func (d *dummyResponseWriter) WriteHeader(i int) {
	d.code = i
}

type serveHTTPTest struct {
	Description   string
	Input         *http.Request
	OutputCode    int
	OutputHeaders http.Header
	OutputBody    string
}

type kv struct {
	key string
	val string
}

func buildRequest(method string, path string, query string, kvx ...kv) *http.Request {
	r := &http.Request{}
	r.Host = "http://localhost"
	r.Header = make(http.Header)
	r.URL = &url.URL{
		Path:     path,
		RawQuery: query,
		Host:     "localhost",
		Scheme:   "http",
	}
	for _, k := range kvx {
		r.Header.Add(k.key, k.val)
	}
	r.Method = method
	return r
}

func buildRequestTLS(method string, path string, query string, kvx ...kv) *http.Request {
	r := &http.Request{}
	r.Host = "https://localhost"
	r.Header = make(http.Header)
	r.URL = &url.URL{
		Path:     path,
		RawQuery: query,
		Host:     "localhost",
		Scheme:   "https",
	}
	r.TLS = &tls.ConnectionState{} // marks the request as TLS
	for _, k := range kvx {
		r.Header.Add(k.key, k.val)
	}
	r.Method = method
	return r
}

var defaultHeaders = http.Header(map[string][]string{
	"Cache-Control": []string{"no-cache"},
	"Pragma":        []string{"no-cache"},
	"Content-Type":  []string{"application/jrd+json"},
})

func plusHeader(h http.Header, kvx ...kv) http.Header {
	var h2 = make(http.Header)
	for k, vx := range h {
		for _, v := range vx {
			h2.Add(k, v)
		}
	}

	for _, k := range kvx {
		h2.Add(k.key, k.val)
	}

	return h2
}

func compareHeaders(h1 http.Header, h2 http.Header) bool {
	if len(h1) != len(h2) {
		return false
	}
	keys := []string{}
	for k := range h1 {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if !reflect.DeepEqual(h1[k], h2[k]) {
			return false
		}
	}
	return true
}

func TestServiceServeHTTP(t *testing.T) {

	var tests = []serveHTTPTest{

		{"GET root URL should return 404 not found with no headers",
			buildRequestTLS("GET", "/", ""), http.StatusNotFound, make(http.Header), ""},
		{"POST root URL should return 404 not found with no headers",
			buildRequestTLS("POST", "/", ""), http.StatusNotFound, make(http.Header), ""},
		{"GET /.well-known should return 404 not found with no headers",
			buildRequestTLS("GET", "/.well-known", ""), http.StatusNotFound, make(http.Header), ""},
		{"POST /.well-known should return 404 not found with no headers",
			buildRequestTLS("POST", "/.well-known", ""), http.StatusNotFound, make(http.Header), ""},

		/*
			4.2.  Performing a WebFinger Query
		*/

		/*
			A WebFinger client issues a query using the GET method to the well-
			known [3] resource identified by the URI whose path component is
			"/.well-known/webfinger" and whose query component MUST include the
			"resource" parameter exactly once and set to the value of the URI for
			which information is being sought.
		*/
		{"GET Webfinger resource should return valid response with default headers",
			buildRequestTLS("GET", WebFingerPath, "resource=acct:hello@domain"), http.StatusOK, defaultHeaders,
			`{"links":[{"href":"x","ref":"x"},{"href":"y","ref":"y"},{"href":"z","ref":"z"}]}`},

		{"POST Webfinger URL should fail with MethodNotAllowed",
			buildRequestTLS("POST", WebFingerPath, ""), http.StatusMethodNotAllowed, defaultHeaders, ""},
		{"GET multiple resources should fail with BadRequest",
			buildRequestTLS("GET", WebFingerPath, "resource=acct:hello@domain&resource=acct:hello2@domain"), http.StatusBadRequest, defaultHeaders, ""},

		/*
		   The "rel" parameter MAY be included multiple times in order to
		   request multiple link relation types.
		*/
		{"GET Webfinger resource with rel should filter results",
			buildRequestTLS("GET", WebFingerPath, "resource=acct:hello@domain&rel=x&rel=y"), http.StatusOK, defaultHeaders,
			`{"links":[{"href":"x","ref":"x"},{"href":"y","ref":"y"}]}`},

		/*
		   A WebFinger resource MAY redirect the client; if it does, the
		   redirection MUST only be to an "https" URI and the client MUST
		   perform certificate validation again when redirected.
		*/
		{"GET non-TLS should redirect to TLS",
			buildRequest("GET", WebFingerPath, "resource=acct:hello@domain"), http.StatusSeeOther, plusHeader(
				defaultHeaders, kv{"Location", "https://localhost/.well-known/webfinger?resource=acct:hello@domain"}), ""},

		/*
				A WebFinger resource MUST return a JRD as the representation for the
			   resource if the client requests no other supported format explicitly
			   via the HTTP "Accept" header.  The client MAY include the "Accept"
			   header to indicate a desired representation; representations other
			   than JRD might be defined in future specifications.  The WebFinger
			   resource MUST silently ignore any requested representations that it
			   does not understand or support.  The media type used for the JSON
			   Resource Descriptor (JRD) is "application/jrd+json" (see Section
			   10.2).
		*/
		{"GET with Accept should return as normal",
			buildRequestTLS("GET", WebFingerPath, "resource=acct:hello@domain", kv{"Accept", "application/json"}), http.StatusOK, defaultHeaders,
			`{"links":[{"href":"x","ref":"x"},{"href":"y","ref":"y"},{"href":"z","ref":"z"}]}`},

		/*
		   If the "resource" parameter is a value for which the server has no
		   information, the server MUST indicate that it was unable to match the
		   request as per Section 10.4.5 of RFC 2616.
		*/
		{"GET with a missing user should return 404",
			buildRequestTLS("GET", WebFingerPath, "resource=acct:missinguser@domain"), http.StatusNotFound, defaultHeaders, ""},

		/*
			If the "resource" parameter is absent or malformed, the WebFinger
			resource MUST indicate that the request is bad as per Section 10.4.1
			of RFC 2616 [2]. (400 bad request)
		*/
		{"GET with no resource should fail with BadRequest",
			buildRequestTLS("GET", WebFingerPath, ""), http.StatusBadRequest, defaultHeaders, ""},
		{"GET with malformed resource URI should fail with BadRequest",
			buildRequestTLS("GET", WebFingerPath, "resource=hello-world"), http.StatusBadRequest, defaultHeaders, ""},
		{"GET with http resource URI should fail with BadRequest",
			buildRequestTLS("GET", WebFingerPath, "resource=http://hello-world"), http.StatusBadRequest, defaultHeaders, ""},
	}
	svc := Default(&dummyUserResolver{})

	for _, tx := range tests {
		w := &dummyResponseWriter{code: 200}
		svc.ServeHTTP(w, tx.Input)
		// code should be 404
		// headers should be empty
		body := strings.TrimSpace(string(w.Buffer.Bytes()))

		failed := false
		failed = failed || w.code != tx.OutputCode
		failed = failed || !compareHeaders(tx.OutputHeaders, w.headers)
		failed = failed || body != tx.OutputBody
		if failed {
			t.Errorf("%s\nHTTP '%v' '%v' => \n\t'%v'\n\t'%v'\n\t'%v';\nexpected \n\t'%v \n\t'%v' \n\t'%v'",
				tx.Description,
				tx.Input.Method, tx.Input.URL,
				w.code, w.headers, body,
				tx.OutputCode, tx.OutputHeaders, tx.OutputBody)
		}
	}

}
