// This file only provides command to generate
// the assets subpackage

//go:generate go-bindata -o assets.go -pkg=assets -prefix=dist/ dist/...

package assets
