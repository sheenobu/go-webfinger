package webfinger

// A Link is a series of user details
type Link struct {
	HRef       string             `json:"href"`
	Type       string             `json:"type,omitempty"`
	Rel        string             `json:"ref"`
	Properties map[string]*string `json:"properties,omitempty"`
	Titles     map[string]string  `json:"titles,omitempty"`
}

// Rel allows referencing a subset of the users details
type Rel string
