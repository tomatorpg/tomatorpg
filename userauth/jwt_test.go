package userauth_test

import (
	"testing"
	"time"

	"github.com/tomatorpg/tomatorpg/userauth"

	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
)

func TestDecodeTokenStr(t *testing.T) {
	key := "abcdef"
	claims := jws.Claims{
		"hello": "world",
		"foo":   "bar",
	}
	jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, _ := jwtToken.Serialize([]byte(key))

	parsedToken, err := userauth.DecodeTokenStr(key, string(serializedToken))
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	for claimName, value := range jwtToken.Claims() {
		if !parsedToken.Claims().Has(claimName) {
			t.Errorf("result token does not have %#v claim", claimName)
		} else if want, have := value, parsedToken.Claims()[claimName]; want != have {
			t.Errorf("expected %#v for claim %#v, got %#v", want, claimName, have)
		}
	}
}

func TestDecodeTokenStr_error(t *testing.T) {
	key := "abcdef"
	claims := jws.Claims{
		"hello": "world",
		"foo":   "bar",
	}
	jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, _ := jwtToken.Serialize([]byte(key))

	_, err := userauth.DecodeTokenStr("wrongkey", string(serializedToken))
	if err == nil {
		t.Errorf("expected error and got nil")
	}
	if want, have := "error validating token: signature is invalid", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestDecodeTokenStr_expired(t *testing.T) {
	key := "abcdef"
	claims := jws.Claims{
		"hello": "world",
		"foo":   "bar",
	}
	claims.SetExpiration(time.Now().Add(-60 * time.Second)) // expired 60 seconds before
	jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, _ := jwtToken.Serialize([]byte(key))

	_, err := userauth.DecodeTokenStr(key, string(serializedToken))
	if err == nil {
		t.Errorf("expected error and got nil")
	}
	if want, have := "error validating token: token is expired", err.Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestEncodeTokenStr(t *testing.T) {
	key := "tyuiop"
	claims := jws.Claims{
		"hello": "world",
		"foo":   "bar",
	}
	tokenStr, err := userauth.EncodeTokenStr(key, claims)
	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	parsedToken, err := userauth.DecodeTokenStr(key, tokenStr)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	for claimName, value := range claims {
		if !parsedToken.Claims().Has(claimName) {
			t.Errorf("result token does not have %#v claim", claimName)
		} else if want, have := value, parsedToken.Claims()[claimName]; want != have {
			t.Errorf("expected %#v for claim %#v, got %#v", want, claimName, have)
		}
	}
}
