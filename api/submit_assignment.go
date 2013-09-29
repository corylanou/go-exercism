package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"corylanou/go-exercism/api/configuration"
)

func SubmitAssignment(config configuration.Config, filePath string, code []byte) (submissionPath string, err error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s", config.Hostname, path)

	submission := struct {
		Key  string `json:"key"`
		Code string `json:"code"`
		Path string `json:"path"`
	}{
		Key:  config.ApiKey,
		Code: string(code),
		Path: filePath,
	}

	var submissionJson []byte
	if submissionJson, err = json.Marshal(submission); err != nil {
		return
	}

	var req *http.Request
	if req, err = http.NewRequest("POST", url, bytes.NewReader(submissionJson)); err != nil {
		return
	}

	req.Header.Set("User-Agent", fmt.Sprintf("github.com/kytrinyx/exercism CLI v%s", Version))

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		err = fmt.Errorf("Error posting assignment: [%s]", err)
		return
	}
	defer resp.Body.Close()

	var response struct {
		Status         string
		Language       string
		Exercise       string
		SubmissionPath string `json:"submission_path"`
		Error          string
	}

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&response); err != nil {
		err = fmt.Errorf("Error parsing API response: [%s]", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("Status: %d, Error: %s", resp.StatusCode, response.Error)
		return
	}

	submissionPath = response.SubmissionPath
	return
}
