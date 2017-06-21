package userauth

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/jose.v1/jws"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GithubConfig provides OAuth2 config for google login
func GithubConfig(hostURL string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  hostURL + "/oauth2/github/callback",
		ClientID:     os.Getenv("OAUTH2_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GITHUB_CLIENT_SECRET"),
		Scopes: []string{
			"user:email",
		},
		Endpoint: github.Endpoint,
	}
}

// GithubCallback returns a http.Handler for Google account login handing
func GithubCallback(conf *oauth2.Config, db *gorm.DB, jwtKey, hostURL string) http.HandlerFunc {
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
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return
		}

		// read into
		/*
			// NOTE: JSON structure of normal response body
			{
			  "login": "octocat",
			  "id": 1,
			  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
			  "gravatar_id": "",
			  "url": "https://api.github.com/users/octocat",
			  "html_url": "https://github.com/octocat",
			  "followers_url": "https://api.github.com/users/octocat/followers",
			  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
			  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
			  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
			  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
			  "organizations_url": "https://api.github.com/users/octocat/orgs",
			  "repos_url": "https://api.github.com/users/octocat/repos",
			  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
			  "received_events_url": "https://api.github.com/users/octocat/received_events",
			  "type": "User",
			  "site_admin": false,
			  "name": "monalisa octocat",
			  "company": "GitHub",
			  "blog": "https://github.com/blog",
			  "location": "San Francisco",
			  "email": "octocat@github.com",
			  "hireable": false,
			  "bio": "There once was...",
			  "public_repos": 2,
			  "public_gists": 1,
			  "followers": 20,
			  "following": 0,
			  "created_at": "2008-01-14T04:33:35Z",
			  "updated_at": "2008-01-14T04:33:35Z",
			  "total_private_repos": 100,
			  "owned_private_repos": 100,
			  "private_gists": 81,
			  "disk_usage": 10000,
			  "collaborators": 8,
			  "two_factor_authentication": true,
			  "plan": {
			    "name": "Medium",
			    "space": 400,
			    "private_repos": 20,
			    "collaborators": 0
			  }
			}
		*/

		result := lzjson.Decode(resp.Body)
		log.Printf("raw result: %s", result.Raw())

		if emailNode := result.Get("email"); emailNode.Type() != lzjson.TypeString || emailNode.String() == "" {
			// TODO: read email from email endpoint
			// TODO: or redirect to error handling page
			log.Printf("no email!")
		}

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
		logger.Log(
			"message", "user found or created",
			"user.id", authUser.ID,
			"user.name", authUser.Name,
		)

		// Create JWS claims with the user info
		expires := time.Now().Add(7 * 24 * time.Hour) // 7 days later
		claims := jws.Claims{}
		claims.Set("id", authUser.ID)
		claims.Set("name", authUser.Name)
		claims.SetAudience("localhost") // TODO: set audience correctly
		claims.SetExpiration(expires)
		tokenStr, _ := EncodeTokenStr(jwtKey, claims)

		http.SetCookie(w, &http.Cookie{
			Name:     "tomatorpg-token",
			Value:    tokenStr,
			Expires:  expires,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, hostURL, http.StatusFound)
	}
}
