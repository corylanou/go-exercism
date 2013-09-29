package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"corylanou/go-exercism/api/configuration"
	"github.com/stretchr/testify/assert"
)

var submitHandler = func(rw http.ResponseWriter, r *http.Request) {
	pathMatches := r.URL.Path == "/api/v1/user/assignments"
	methodMatches := r.Method == "POST"
	if !(pathMatches && methodMatches) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	userAgentMatches := r.Header.Get("User-Agent") == fmt.Sprintf("github.com/kytrinyx/exercism CLI v%s", Version)

	if !userAgentMatches {
		fmt.Printf("User agent mismatch: %s\n", r.Header.Get("User-Agent"))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Reading body error: %s\n", err)
		return
	}

	type Submission struct {
		Key  string
		Code string
		Path string
	}

	submission := Submission{}

	err = json.Unmarshal(body, &submission)
	if err != nil {
		fmt.Printf("Unmarshalling error: %s, Body: %s\n", err, body)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if submission.Key != "myApiKey" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(rw, `{"error": "Unable to identify user"}`)
		return
	}

	code := submission.Code
	filePath := submission.Path

	codeMatches := string(code) == "My source code\n"
	filePathMatches := filePath == "ruby/bob/bob.rb"

	if !filePathMatches {
		fmt.Printf("FilePathMismatch: File Path: %s\n", filePath)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !codeMatches {
		fmt.Printf("Code Mismatch: Code: %v\n", code)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	submitJson := `
{
	"status":"saved",
	"language":"ruby",
	"exercise":"bob",
	"submission_path":"/username/ruby/bob"
}
`
	fmt.Fprintf(rw, submitJson)
}

func TestSubmitWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))
	defer server.Close()

	var code = []byte("My source code\n")
	config := configuration.Config{
		Hostname: server.URL,
		ApiKey:   "myApiKey",
	}
	response, err := SubmitAssignment(config, "ruby/bob/bob.rb", code)
	assert.NoError(t, err)

	// We don't use these values in any returns that I can find, only to talk to the server.
	// Is this the intent or do we plan on using them?  If not, I prefer to keep the data local.
	//assert.Equal(t, response.Status, "saved")
	//assert.Equal(t, response.Language, "ruby")
	//assert.Equal(t, response.Exercise, "bob")
	//assert.Equal(t, response.SubmissionPath, "/username/ruby/bob")
	assert.Equal(t, response, "/username/ruby/bob")
}

func TestSubmitWithIncorrectKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))
	defer server.Close()

	config := configuration.Config{
		Hostname: server.URL,
		ApiKey:   "myWrongApiKey",
	}

	var code = []byte("My source code\n")
	response, err := SubmitAssignment(config, "ruby/bob/bob.rb", code)

	assert.Error(t, err)
	assert.Empty(t, response)
}
