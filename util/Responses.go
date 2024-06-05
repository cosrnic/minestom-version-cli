package util

type GHSuccessResponse struct {
	Sha     string   `json:"sha"`
	Commit  Commit   `json:"commit"`
	Parents []Parent `json:"parents"`
}

type Parent struct {
	Sha string `json:"sha"`
}
type Commit struct {
	Author  Author `json:"author"`
	Message string `json:"message"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type GHErrorResponse struct {
	Message string `json:"message"`
}

type GHCheckRunsSuccessResponse struct {
	CheckRuns []CheckRuns `json:"check_runs"`
}

type CheckRuns struct {
	Conclusion string `json:"conclusion"`
}

type RateLimit struct {
	Message string `json:"message"`
}
