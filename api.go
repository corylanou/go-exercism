package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const VERSION = "1.0.1"

var FetchEndpoints = map[string]string{
	"current": "/api/v1/user/assignments/current",
	"next":    "/api/v1/user/assignments/next",
	"demo":    "/api/v1/assignments/demo",
}

type submitResponse struct {
	Status         string
	Language       string
	Exercise       string
	SubmissionPath string `json:"submission_path"`
}

func FetchAssignments(config Config, path string) (as []Assignment, err error) {
	url := fmt.Sprintf("%s%s?key=%s", config.Hostname, path, config.ApiKey)

	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		err = fmt.Errorf("Error fetching assignments: [%s]", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error fetching assignments. HTTP Status Code: %d", resp.StatusCode)
		return
	}

	var response struct {
		Assignments []Assignment
	}

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&response); err != nil {
		err = fmt.Errorf("Error parsing API response: [%s]", err)
		return
	}

	as = response.Assignments
	return
}

func SubmitAssignment(config Config, filePath string, code []byte) (r *submitResponse, err error) {
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

	req.Header.Set("User-Agent", fmt.Sprintf("github.com/kytrinyx/exercism CLI v%s", VERSION))

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		err = fmt.Errorf("Error posting assignment: [%s]", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		response := struct {
			Error string
		}{}
		dec := json.NewDecoder(resp.Body)
		if err = dec.Decode(&response); err != nil {
			err = fmt.Errorf("Error parsing API response: [%s]", err)
			return
		}
		err = fmt.Errorf("Status: %d, Error: %s", resp.StatusCode, response.Error)
		return
	}

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&r); err != nil {
		err = fmt.Errorf("Error parsing API response: [%s]", err)
		return
	}

	return
}
