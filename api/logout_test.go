package api

import (
	"io/ioutil"
	"os"
	"testing"

	"corylanou/go-exercism/api/configuration"
	"github.com/stretchr/testify/assert"
)

func asserFileDoesNotExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err == nil {
		t.Errorf("File [%s] already exist.", filename)
	}
}

func TestLogoutDeletesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	c := configuration.Config{}

	configuration.ToFile(tmpDir, c)

	Logout(tmpDir)

	asserFileDoesNotExist(t, configuration.Filename(tmpDir))
}
