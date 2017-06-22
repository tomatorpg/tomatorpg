package userauth

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GithubConfig provides OAuth2 config for google login
func GithubConfig(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("OAUTH2_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_GITHUB_CLIENT_SECRET"),
		Scopes: []string{
			"user:email",
		},
		Endpoint: github.Endpoint,
	}
}

// GithubCallback returns a http.Handler for Google account login handing
func GithubCallback(
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

		// read into
		var primaryEmail string
		var verifiedEmails []string

		client := conf.Client(context.Background(), token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return
		}
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
		userInfoResult := lzjson.Decode(resp.Body)

		// read other email(s) from email endpoint
		// TODO: or redirect to error handling page
		resp, err = client.Get("https://api.github.com/user/emails")
		if err != nil {
			return
		}
		/*
			// NOTE: JSON structure of normal response body
			[
			  {
			    "email": "octocat@github.com",
			    "verified": true,
			    "primary": true,
			    "visibility": "public"
			  }
			]
		*/
		userEmailResult := lzjson.Decode(resp.Body)
		emails := []struct {
			Email    string `json:"email"`
			Verified bool   `json:"verified"`
			Primary  bool   `json:"primary"`
		}{}
		if err = userEmailResult.Unmarshal(&emails); err == nil {
			verifiedEmails = make([]string, 0, len(emails))
			for _, email := range emails {
				if email.Primary {
					primaryEmail = email.Email
				} else if email.Verified {
					verifiedEmails = append(verifiedEmails, email.Email)
				}
			}
		} else {
			log.Printf("unexpected error: %s", err.Error())
			logger.Log(
				"message", "error reading results from github's user/emails endpoint",
				"error", err.Error(),
				"rawResponse", string(userEmailResult.Raw()),
			)
		}

		// TODO: detect read  / decode error
		// TODO: check if the email has been verified or not
		authUser, err := loadOrCreateUser(
			db,
			models.User{
				Name:         userInfoResult.Get("name").String(),
				PrimaryEmail: primaryEmail,
			},
			verifiedEmails,
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
