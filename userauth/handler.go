package userauth

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mrjones/oauth"
	"github.com/tomatorpg/tomatorpg/utils"

	"golang.org/x/oauth2"
)

// AuthURLFactory manufactures redirectURLs to authentication endpoint
// with the correct callback path back to the application site.
type AuthURLFactory func(r *http.Request) (redirectURL string, err error)

// OAuth2AuthURLFactory generates factory of authentication URL
// to the oauth2 config
func OAuth2AuthURLFactory(conf *oauth2.Config) AuthURLFactory {
	return func(r *http.Request) (url string, err error) {
		url = conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		return
	}
}

// OAuth1aConsumer provide ways to get reuest token and auth url
type OAuth1aConsumer interface {
	GetRequestTokenAndUrl(callbackURL string) (token *oauth.RequestToken, url string, err error)
}

// CookieFactory generates cookie struct for auth interactions
// (i.e. login session and logout)
type CookieFactory func(r *http.Request) *http.Cookie

// OAuth1aAuthURLFactory generates factory of authentication URL
// to the oauth1a consumer and callback URL
func OAuth1aAuthURLFactory(c OAuth1aConsumer, callbackURL string) AuthURLFactory {
	return func(r *http.Request) (url string, err error) {
		logger := utils.GetLogger(r.Context())
		requestToken, url, err := c.GetRequestTokenAndUrl(callbackURL)
		if err != nil {
			logger.Log(
				"at", "error",
				"message", "error retrieving twitter token",
				"error", err.Error(),
			)
			return
		}
		TokenSave(requestToken)
		return
	}
}

// RedirectHandler handles the generation and redirection to
// authentication endpoint with proper parameters
func RedirectHandler(getAuthURL AuthURLFactory, errURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := getAuthURL(r)
		if err != nil {
			// TODO: redirect to the errURL with status messages
			http.Redirect(w, r, errURL, http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// LogoutHandler makes a cookie of a given name expires
func LogoutHandler(redirectURL string, getLoginCookie CookieFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := getLoginCookie(r)
		cookie.Expires = time.Now().Add(-1 * time.Hour) // expires immediately
		http.SetCookie(w, cookie)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

// LoginHandler return a mux to handle all login related routes
func LoginHandler(
	db *gorm.DB,
	genLoginCookie CookieFactory,
	providers []AuthProvider,
	jwtSecret, baseURL, oauth2Path, successPath, errPath string,
) http.Handler {

	// Note: oauth2Path must start with "/" and must not have trailing slash
	// Note: baseURL must be full URL without path or any trailing slash

	mux := http.NewServeMux()
	oauth2URL := baseURL + oauth2Path   // full URL to oauth2 path
	successURL := baseURL + successPath // full URL to success page

	if provider := FindProvider("google", providers); provider != nil {
		mux.Handle(oauth2Path+"/google", RedirectHandler(
			OAuth2AuthURLFactory(GoogleConfig(
				*provider,
				oauth2URL+"/google/callback",
			)),
			errPath,
		))
		mux.Handle(oauth2Path+"/google/callback", GoogleCallback(
			GoogleConfig(
				*provider,
				oauth2URL+"/google/callback",
			),
			db,
			genLoginCookie,
			jwtSecret,
			successURL,
			errPath,
		))
	}

	if provider := FindProvider("facebook", providers); provider != nil {
		mux.Handle(oauth2Path+"/facebook", RedirectHandler(
			OAuth2AuthURLFactory(
				FacebookConfig(
					*provider,
					oauth2URL+"/facebook/callback",
				),
			),
			errPath,
		))
		mux.Handle(oauth2Path+"/facebook/callback", FacebookCallback(
			FacebookConfig(
				*provider,
				oauth2URL+"/facebook/callback",
			),
			db,
			genLoginCookie,
			jwtSecret,
			successURL,
			errPath,
		))
	}

	if provider := FindProvider("github", providers); provider != nil {
		mux.Handle(oauth2Path+"/github", RedirectHandler(
			OAuth2AuthURLFactory(GithubConfig(
				*provider,
				oauth2URL+"/github/callback",
			)),
			errPath,
		))
		mux.Handle(oauth2Path+"/github/callback", GithubCallback(
			GithubConfig(
				*provider,
				oauth2URL+"/github/callback",
			),
			db,
			genLoginCookie,
			jwtSecret,
			successURL,
			errPath,
		))
	}

	if provider := FindProvider("twitter", providers); provider != nil {
		mux.Handle(oauth2Path+"/twitter", RedirectHandler(
			OAuth1aAuthURLFactory(
				TwitterConsumer(*provider),
				oauth2URL+"/twitter/callback",
			),
			errPath,
		))
		mux.Handle(oauth2Path+"/twitter/callback", TwitterCallback(
			TwitterConsumer(*provider),
			db,
			TokenConsume,
			genLoginCookie,
			jwtSecret,
			successURL,
			errPath,
		))
	}

	return mux
}
