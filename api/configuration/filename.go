package configuration

import (
	"fmt"
)

const fileExt = ".exercism.go"

func Filename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, fileExt)
}
