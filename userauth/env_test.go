package userauth_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tomatorpg/tomatorpg/userauth"
)

func genGetEnv(id string) func(string) string {
	clientIDKey := fmt.Sprintf("OAUTH2_%s_CLIENT_ID", strings.ToUpper(id))
	clientSecretKey := fmt.Sprintf("OAUTH2_%s_CLIENT_SECRET", strings.ToUpper(id))
	return func(key string) (value string) {
		if key == clientIDKey || key == clientSecretKey {
			value = key
		}
		return
	}
}

func chainGetEnv(genEnvFns ...func(string) string) func(string) string {
	return func(key string) (value string) {
		for _, getEnv := range genEnvFns {
			if value = getEnv(key); value != "" {
				return
			}
		}
		return
	}
}

func mapProviderID(providers []userauth.AuthProvider) (ids []string) {
	ids = make([]string, len(providers))
	for i, provider := range providers {
		ids[i] = provider.ID
	}
	return
}

func checkProvider(id string) func([]userauth.AuthProvider) error {
	return func(providers []userauth.AuthProvider) (err error) {
		for _, provider := range providers {
			if provider.ID == id {
				return
			}
		}
		err = fmt.Errorf("provider %s not found", id)
		return
	}
}

func compare(arr1 []string, arr2 []string) (err error) {
	for _, str1 := range arr1 {
		str1InArr2 := false
		for _, str2 := range arr2 {
			if str1 == str2 {
				str1InArr2 = true
				break
			}
		}
		if !str1InArr2 {
			err = fmt.Errorf("does not have %s", str1)
		}
	}
	return
}

func TestEnvProviders(t *testing.T) {
	tests := []struct {
		getEnv func(string) string
		ids    []string
	}{
		{
			getEnv: chainGetEnv(),
			ids:    []string{},
		},
		{
			getEnv: genGetEnv("google"),
			ids:    []string{"google"},
		},
		{
			getEnv: genGetEnv("facebook"),
			ids:    []string{"facebook"},
		},
		{
			getEnv: genGetEnv("twitter"),
			ids:    []string{"twitter"},
		},
		{
			getEnv: genGetEnv("github"),
			ids:    []string{"github"},
		},
		{
			getEnv: chainGetEnv(genGetEnv("google"), genGetEnv("github")),
			ids:    []string{"google", "github"},
		},
	}
	for _, test := range tests {
		providers := userauth.EnvProviders(test.getEnv, "/")
		providerIDs := mapProviderID(providers)
		if err := compare(test.ids, providerIDs); err != nil {
			t.Errorf("providerIDs missed ID: %s", err.Error())
		}
		if err := compare(providerIDs, test.ids); err != nil {
			t.Errorf("providerIDs has more ID than expected: %s", err.Error())
		}
	}
}

func TestFindProvider(t *testing.T) {
	found := userauth.FindProvider(
		"world",
		[]userauth.AuthProvider{
			{
				ID: "hello",
			},
			{
				ID: "world",
			},
		},
	)
	if want, have := "world", found.ID; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	found = userauth.FindProvider(
		"foo",
		[]userauth.AuthProvider{
			{
				ID: "hello",
			},
			{
				ID: "world",
			},
		},
	)
	if found != nil {
		t.Errorf("expected nil, got %#v", found)
	}
}
