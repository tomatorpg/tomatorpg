package pubsub_test

import (
	"testing"

	"github.com/tomatorpg/tomatorpg/protocol/pubsub"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
)

func TestParseJWT(t *testing.T) {
	key := "abcdef"
	claims := jws.Claims{
		"hello": "world",
		"foo":   "bar",
	}
	jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, _ := jwtToken.Serialize([]byte(key))

	parsedToken, err := pubsub.ParseToken(key, string(serializedToken))
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
