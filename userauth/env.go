package userauth

import (
	"fmt"
	"strings"
)

// AuthProvider defines a login provider in details
type AuthProvider struct {
	ID           string
	Name         string
	Path         string
	ClientID     string
	ClientSecret string
}

// EnvProviders gets login providers from environment
func EnvProviders(getEnv func(string) string, basePath string) (providers []AuthProvider) {
	providers = make([]AuthProvider, 0, 4)
	protoProviders := []AuthProvider{
		{
			ID:   "google",
			Name: "Google",
			Path: basePath + "/google",
		},
		{
			ID:   "facebook",
			Name: "Facebook",
			Path: basePath + "/facebook",
		},
		{
			ID:   "twitter",
			Name: "Twitter",
			Path: basePath + "/twitter",
		},
		{
			ID:   "github",
			Name: "Github",
			Path: basePath + "/github",
		},
	}

	// read client id and key from environment
	var clientIDKey, clientSecretKey string
	for _, provider := range protoProviders {
		clientIDKey = fmt.Sprintf("OAUTH2_%s_CLIENT_ID", strings.ToUpper(provider.ID))
		clientSecretKey = fmt.Sprintf("OAUTH2_%s_CLIENT_SECRET", strings.ToUpper(provider.ID))
		if clientID, clientSecret := getEnv(clientIDKey), getEnv(clientSecretKey); clientID != "" && clientSecret != "" {
			provider.ClientID, provider.ClientSecret = clientID, clientSecret
			providers = append(
				providers,
				provider,
			)
		}
	}
	return
}

// FindProvider find provider of given ID, or return nil
func FindProvider(id string, providers []AuthProvider) *AuthProvider {
	for _, provider := range providers {
		if provider.ID == id {
			return &provider
		}
	}
	return nil
}
