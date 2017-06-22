package userauth

import (
	"log"
	"net/http"
	"os"

	"github.com/go-restit/lzjson"
	"github.com/jinzhu/gorm"
	"github.com/mrjones/oauth"
	"github.com/tomatorpg/tomatorpg/models"
	"github.com/tomatorpg/tomatorpg/utils"
)

var tokens map[string]*oauth.RequestToken

func init() {
	tokens = make(map[string]*oauth.RequestToken, 1024)
}

// TokenSave stores a copy of token in a map by token key
func TokenSave(token *oauth.RequestToken) {
	// TODO: add mutex lock mechanism
	tokens[token.Token] = token
}

// TokenConsume return the token stored previously and remove it
// from the map
func TokenConsume(tokenKey string) (token *oauth.RequestToken) {
	// TODO: add mutex lock mechanism
	token, ok := tokens[tokenKey]
	if ok {
		delete(tokens, tokenKey)
	}
	return
}

// TwitterConsumer provides OAuth config for twitter login
func TwitterConsumer() *oauth.Consumer {
	return oauth.NewConsumer(
		os.Getenv("OAUTH2_TWITTER_CLIENT_ID"),
		os.Getenv("OAUTH2_TWITTER_CLIENT_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
}

// TwitterCallback returns a http.Handler for Twitter account login handing
func TwitterCallback(
	c *oauth.Consumer,
	db *gorm.DB,
	tokenConsume func(tokenKey string) *oauth.RequestToken,
	genLoginCookie CookieFactory,
	jwtKey, hostURL string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := utils.GetLogger(r.Context())
		values := r.URL.Query()
		verificationCode := values.Get("oauth_verifier")
		tokenKey := values.Get("oauth_token")

		accessToken, err := c.AuthorizeToken(tokenConsume(tokenKey), verificationCode)
		if err != nil {
			log.Fatal(err)
		}

		client, err := c.MakeHttpClient(accessToken)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Get(
			"https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true&skip_status")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// read into
		result := lzjson.Decode(resp.Body)
		/*
			// NOTE: JSON structure of normal response body
			{
				"contributors_enabled": true,
				"created_at": "Sat May 09 17:58:22 +0000 2009",
				"default_profile": false,
				"default_profile_image": false,
				"description": "I taught your phone that thing you like.  The Mobile Partner Engineer @Twitter. ",
				"favourites_count": 588,
				"follow_request_sent": null,
				"followers_count": 10625,
				"following": null,
				"friends_count": 1181,
				"geo_enabled": true,
				"id": 38895958,
				"id_str": "38895958",
				"is_translator": false,
				"lang": "en",
				"listed_count": 190,
				"location": "San Francisco",
				"name": "Sean Cook",
				"email": "sean.cook@email.com",
				"notifications": null,
				"profile_background_color": "1A1B1F",
				"profile_background_image_url": "http://a0.twimg.com/profile_background_images/495742332/purty_wood.png",
				"profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/495742332/purty_wood.png",
				"profile_background_tile": true,
				"profile_image_url": "http://a0.twimg.com/profile_images/1751506047/dead_sexy_normal.JPG",
				"profile_image_url_https": "https://si0.twimg.com/profile_images/1751506047/dead_sexy_normal.JPG",
				"profile_link_color": "2FC2EF",
				"profile_sidebar_border_color": "181A1E",
				"profile_sidebar_fill_color": "252429",
				"profile_text_color": "666666",s
				"profile_use_background_image": true,
				"protected": false,
				"screen_name": "theSeanCook",
				"show_all_inline_media": true,
				"statuses_count": 2609,
				"time_zone": "Pacific Time (US & Canada)",
				"url": null,
				"utc_offset": -28800,
				"verified": false
			}
		*/

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
