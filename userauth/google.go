package userauth

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig(hostURL string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  hostURL + "/oauth2/google/callback",
		ClientID:     os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleCallback(conf *oauth2.Config, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Printf("Code exchange failed with '%s'\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		client := conf.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			return
		}

		// read into
		/*
			// NOTE: JSON structure of normal response body
			{
			 "id": "some-id-in-google-account",
			 "email": "email-for-the-account",
			 "verified_email": true,
			 "name": "Some Name",
			 "given_name": "Some",
			 "family_name": "Name",
			 "link": "https://plus.google.com/+SomeUserOnGPlus",
			 "picture": "url-to-some-picture",
			 "gender": "female",
			 "locale": "zh-HK"
			}
		*/

		result := lzjson.Decode(resp.Body)
		// TODO: detect read  / decode error
		// TODO: check if the email has been verified or not
		authUser := models.User{
			Name:         result.Get("name").String(),
			PrimaryEmail: result.Get("email").String(),
		}

		// search existing user with the email
		var userEmail models.UserEmail
		var prevUser models.User

		if db.First(&prevUser, "primary_email = ?", authUser.PrimaryEmail); prevUser.PrimaryEmail != "" {
			// TODO: log this?
			authUser = prevUser
		} else if db.First(&userEmail, "email = ?", authUser.PrimaryEmail); userEmail.Email != "" {
			// TODO: log this?
			db.First(&authUser, "id = ?", userEmail.UserID)
		} else {

			tx := db.Begin()

			// create user
			if res := tx.Create(&authUser); res.Error != nil {
				// TODO: log and provide error to user
				tx.Rollback()
				return
			}

			// create user-email relation
			newUserEmail := models.UserEmail{
				UserID: authUser.ID,
				Email:  authUser.PrimaryEmail,
			}
			if res := tx.Create(&newUserEmail); res.Error != nil {
				tx.Rollback()
				return
			}

			tx.Commit()
		}
		log.Printf("user found or created: %#v", authUser)

		// Create JWS claims with the user info
		claims := jws.Claims{}
		claims.Set("id", authUser.ID)
		claims.Set("name", authUser.Name)
		claims.SetAudience("localhost")

		jwtToken := jws.NewJWT(claims, crypto.SigningMethodHS256)
		serializedToken, _ := jwtToken.Serialize([]byte("abcdef"))

		http.SetCookie(w, &http.Cookie{
			Name:     "tomatorpg-token",
			Value:    string(serializedToken),
			Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
	}
}
