package userauth

// AuthProvider defines a login provider in details
type AuthProvider struct {
	ID   string
	Name string
	Path string
}

// EnvProviders gets login providers from environment
func EnvProviders(basePath string) []AuthProvider {
	// TODO: render this list based on environment variables
	return []AuthProvider{
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
}
