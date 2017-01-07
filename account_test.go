package webfinger

import (
	"testing"

	"reflect"

	"github.com/pkg/errors"
)

type acp struct {
	Input  string
	Output account
	Error  string
}

func (a *acp) Invoke() error {
	var ax account
	err := ax.ParseString(a.Input)

	failed := false
	failed = failed || err != nil && a.Error == ""
	failed = failed || err == nil && a.Error != ""
	failed = failed || err != nil && a.Error != err.Error()
	failed = failed || err != nil && a.Error != errors.Cause(err).Error()

	failed = failed || !reflect.DeepEqual(a.Output, ax)

	if failed {
		return errors.Errorf("ax.ParseString('%v') => '%v', '%v'; expected '%v', '%v'", a.Input, err, ax, a.Error, a.Output)
	}
	return nil
}

func TestAccountParse(t *testing.T) {
	var tests = []acp{
		{"", account{}, "URI is not an account"},
		{"http://hello.world", account{}, "URI is not an account"},

		{"acct:hello", account{"hello", ""}, "No domain on account"},
		{"acct:hello@domain", account{"hello", "domain"}, ""},

		{"acct:hello@domain/uri", account{"hello", "domain"}, ""},
	}

	for _, tx := range tests {
		if err := tx.Invoke(); err != nil {
			t.Error(err.Error())
		}
	}
}
