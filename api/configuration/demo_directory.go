package configuration

import (
	"fmt"
	"os"
)

func DemoDirectory() (string, error) {
	dir, err := os.Getwd()
	return fmt.Sprintf("%s/exercism-demo", dir), err
}
