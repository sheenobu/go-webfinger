package webfinger

import (
	"encoding/json"
	"errors"
	"net/http"
)

// WebFingerPath defines the default path of the webfinger handler.
const WebFingerPath = "/.well-known/webfinger"

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//TODO: support host-meta as a path

	path := r.URL.Path
	switch path {
	case WebFingerPath:
		s.Webfinger(w, r)
	default:
		s.NotFoundHandler.ServeHTTP(w, r)
	}
}

// Webfinger is the webfinger handler
func (s *Service) Webfinger(w http.ResponseWriter, r *http.Request) {
	s.runPrehandlers(w, r)

	if r.TLS == nil && s.NoTLSHandler != nil {
		s.NoTLSHandler.ServeHTTP(w, r)
		return
	}

	//NOTE: should this run before or after the pre-run handlers?
	if r.Method != "GET" {
		s.MethodNotSupportedHandler.ServeHTTP(w, r)
		return
	}

	if len(r.URL.Query()["resource"]) != 1 {
		s.MalformedRequestHandler.ServeHTTP(w, addError(r, errors.New("Malformed resource parameter")))
		return
	}
	resource := r.URL.Query().Get("resource")
	var a account
	if err := a.ParseString(resource); err != nil {
		s.MalformedRequestHandler.ServeHTTP(w, addError(r, err))
		return
	}

	relStrings := r.URL.Query()["rel"]
	var rels []Rel
	for _, r := range relStrings {
		rels = append(rels, Rel(r))
	}

	rsc, err := s.Resolver.FindUser(a.Name, a.Hostname, rels)
	if err != nil {
		if !s.Resolver.IsNotFoundError(err) {
			s.ErrorHandler.ServeHTTP(w, addError(r, err))
			return
		}

		rsc, err = s.Resolver.DummyUser(a.Name, a.Hostname, rels)
		if err != nil && !s.Resolver.IsNotFoundError(err) {
			s.ErrorHandler.ServeHTTP(w, addError(r, err))
			return
		} else if s.Resolver.IsNotFoundError(err) {
			s.NotFoundHandler.ServeHTTP(w, r)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(&rsc); err != nil {
		s.ErrorHandler.ServeHTTP(w, addError(r, err))
		return
	}
}

func (s *Service) runPrehandlers(w http.ResponseWriter, r *http.Request) {
	if s.PreHandlers == nil {
		return
	}

	for _, val := range s.PreHandlers {
		if val != nil {
			val.ServeHTTP(w, r)
		}
	}
}

func (s *Service) defaultErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Service) defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (s *Service) defaultMethodNotSupportedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) defaultMalformedRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) defaultNoTLSHandler(w http.ResponseWriter, r *http.Request) {
	u := *r.URL
	u.Scheme = "https"
	w.Header().Set("Location", u.String())
	w.WriteHeader(http.StatusSeeOther)
}
