package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"corylanou/go-exercism/api/configuration"
)

func FetchAssignments(config configuration.Config, path string) (as []Assignment, err error) {
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
