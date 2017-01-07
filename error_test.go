package webfinger

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestAddAndGetError(t *testing.T) {

	{ // test with no error
		ctx := context.Background()
		r := &http.Request{}
		r = r.WithContext(ctx)

		err2 := ErrorFromContext(r.Context())
		if err2 != nil {
			t.Errorf("No error should result in no error")
		}
	}

	{ // Test with addError(nil)
		ctx := context.Background()
		r := &http.Request{}
		r = r.WithContext(ctx)
		r = addError(r, nil)

		err2 := ErrorFromContext(r.Context())
		if err2 != nil {
			t.Errorf("Error was nil, is now not nil")
		}
	}

	{ // Test with addError(errors.New("X"))
		ctx := context.Background()
		r := &http.Request{}
		r = r.WithContext(ctx)
		r = addError(r, errors.New("X"))

		err2 := ErrorFromContext(r.Context())
		if err2 == nil || err2.Error() != "X" {
			t.Errorf("Err is %v, expected 'X'", err2)
		}
	}

}
