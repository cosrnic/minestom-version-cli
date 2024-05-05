package util

type GHSuccessResponse struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author  Author `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type GHErrorResponse struct {
	Message string `json:"message"`
}
