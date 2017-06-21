package userauth_test

import (
	"fmt"
	"testing"

	"github.com/tomatorpg/tomatorpg/userauth"
)

func TestLoginErrorType(t *testing.T) {
	tests := []struct {
		err error
		str string
	}{
		{
			err: userauth.ErrUnknown,
			str: "unknown error",
		},
		{
			err: userauth.ErrNoEmail,
			str: "no email",
		},
		{
			err: userauth.ErrDatabase,
			str: "database error",
		},
	}

	for _, test := range tests {
		if want, have := fmt.Sprintf("LoginError(%#v)", test.str), fmt.Sprintf("%#v", test.err); want != have {
			t.Errorf("expected %#v String() to give %#v, got %#v", test.err, want, have)
		}
		if want, have := "login error: "+test.str, fmt.Sprintf("%s", test.err); want != have {
			t.Errorf("expected %#v String() to give %#v, got %#v", test.err, want, have)
		}
	}
}

func TestLoginError(t *testing.T) {
	tests := []struct {
		err   userauth.LoginError
		str   string
		gostr string
	}{
		{
			err:   userauth.LoginError{},
			str:   "unknown error",
			gostr: `"unknown error"`,
		},
		{
			err: userauth.LoginError{
				Type: userauth.ErrNoEmail,
			},
			str:   "no email",
			gostr: `"no email"`,
		},
		{
			err: userauth.LoginError{
				Type: userauth.ErrDatabase,
			},
			str:   "database error",
			gostr: `"database error"`,
		},
		{
			err: userauth.LoginError{
				Type: userauth.ErrDatabase,
				Err:  fmt.Errorf("some dummy error"),
			},
			str:   `error="some dummy error"`,
			gostr: `error="some dummy error"`,
		},
		{
			err: userauth.LoginError{
				Type:   userauth.ErrDatabase,
				Action: "create dummy",
				Err:    fmt.Errorf("some dummy error"),
			},
			str:   `action="create dummy" error="some dummy error"`,
			gostr: `action="create dummy" error="some dummy error"`,
		},
	}

	for _, test := range tests {
		if want, have := fmt.Sprintf("LoginError(%s)", test.gostr), fmt.Sprintf("%#v", test.err); want != have {
			t.Errorf("expected %#v String() to give %#v, got %#v", test.err, want, have)
		}
		if want, have := "login error: "+test.str, fmt.Sprintf("%s", test.err); want != have {
			t.Errorf("expected %#v String() to give %#v, got %#v", test.err, want, have)
		}
	}

}
