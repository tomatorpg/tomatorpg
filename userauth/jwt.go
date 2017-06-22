package userauth

import (
	"fmt"
	"net/http"

	"github.com/tomatorpg/tomatorpg/models"

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

func authJWTCookie(cookie *http.Cookie, jwtKey string, authUser models.User) *http.Cookie {

	// Create JWS claims with the user info
	claims := jws.Claims{}
	claims.Set("id", authUser.ID)
	claims.Set("name", authUser.Name)
	claims.SetAudience("localhost") // TODO: set audience correctly
	claims.SetExpiration(cookie.Expires)

	// encode token and store in cookies
	cookie.Value, _ = EncodeTokenStr(jwtKey, claims)
	return cookie
}
