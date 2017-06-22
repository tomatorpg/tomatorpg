package userauth

import (
	"context"
	"net/http"
	"os"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// FacebookConfig provides OAuth2 config for google login
func FacebookConfig(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("OAUTH2_FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_FACEBOOK_CLIENT_SECRET"),
		Scopes: []string{
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
}

// FacebookCallback returns a http.Handler for Google account login handing
func FacebookCallback(
	conf *oauth2.Config,
	db *gorm.DB,
	genLoginCookie CookieFactory,
	jwtKey, hostURL string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := utils.GetLogger(r.Context())
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			logger.Log(
				"at", "error",
				"message", "code exchange failed",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		client := conf.Client(context.Background(), token)
		resp, err := client.Get("https://graph.facebook.com/v2.9/me?fields=id,name,email")
		if err != nil {
			return
		}

		// read into
		/*
			// NOTE: JSON structure of normal response body
			{
			  "id": "numerical-user-id",
			  "name": "user display name",
			  "email": "email address"
			}
		*/

		result := lzjson.Decode(resp.Body)

		// TODO: detect read  / decode error
		// TODO: check if the email has been verified or not
		authUser, err := loadOrCreateUser(
			db,
			models.User{
				Name:         result.Get("name").String(),
				PrimaryEmail: result.Get("email").String(),
			},
			[]string{},
		)

		if err != nil {
			logger.Log(
				"at", "error",
				"message", "error",
				"error", err.Error(),
			)
			// TODO; return some warning message to redirected page
			http.Redirect(w, r, hostURL, http.StatusFound)
			return
		}

		logger.Log(
			"at", "info",
			"message", "user found or created",
			"user.id", authUser.ID,
			"user.name", authUser.Name,
		)

		// set authUser digest to cookie as jwt
		http.SetCookie(w,
			authJWTCookie(genLoginCookie(r), jwtKey, *authUser))

		http.Redirect(w, r, hostURL, http.StatusTemporaryRedirect)
	}
}
