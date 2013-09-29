package api

import (
	"os"

	"corylanou/go-exercism/api/configuration"
)

func Logout(dir string) {
	os.Remove(configuration.Filename(dir))
}
