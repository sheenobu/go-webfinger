package webfinger

import (
	"errors"
	"strings"
)

type account struct {
	Name     string
	Hostname string
}

func (a *account) ParseString(str string) (err error) {
	if !strings.HasPrefix(str, "acct:") {
		err = errors.New("URI is not an account")
		return
	}

	items := strings.Split(str, "@")
	a.Name = items[0][5:]
	if len(items) < 2 {
		//TODO: this might not be required
		err = errors.New("No domain on account")
		return
	}

	a.Hostname = strings.Split(items[1], "/")[0]

	return
}
