package webfinger

// A Resource is the top-level JRD resource object.
type Resource struct {
	Subject    string            `json:"subject,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Links      []Link            `json:"links"`
}
