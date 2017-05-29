package pubsub

import (
	"fmt"

	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"gopkg.in/jose.v1/jwt"
)

// ParseToken parses token string into JWT token
func ParseToken(key, tokenStr string) (token jwt.JWT, err error) {
	token, _ = jws.ParseJWT([]byte(tokenStr))
	if err = token.Validate([]byte(key), crypto.SigningMethodHS256); err != nil {
		err = fmt.Errorf("error validating token: %s", err.Error())
		return
	}

	// TODO: further validate token (e.g. expires)
	return
}
