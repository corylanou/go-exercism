package api

type Assignment struct {
	Track    string
	Slug     string
	Readme   string
	TestFile string `json:"test_file"`
	Tests    string
}
