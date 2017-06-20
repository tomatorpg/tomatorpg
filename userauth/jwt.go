package userauth

import (
	"fmt"

	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"gopkg.in/jose.v1/jwt"
)

// DecodeTokenStr parses token string into JWT token
func DecodeTokenStr(key, tokenStr string) (token jwt.JWT, err error) {
	token, _ = jws.ParseJWT([]byte(tokenStr))
	if err = token.Validate([]byte(key), crypto.SigningMethodHS256); err != nil {
		err = fmt.Errorf("error validating token: %s", err.Error())
		return
	}
	return
}

// EncodeTokenStr encode a given claim as a JWT token string
func EncodeTokenStr(key string, claims jws.Claims) (tokenStr string, err error) {
	jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
	serializedToken, err := jwtToken.Serialize([]byte(key))
	tokenStr = string(serializedToken)
	return
}
