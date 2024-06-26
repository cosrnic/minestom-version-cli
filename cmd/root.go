package cmd

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/cosrnic/minestom-version-cli/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "minestom-version",
	Short: "Minestom-version CLI in Go",
	Long:  "A CLI for getting the latest version of Minestom [or a branch of minestom]",
	Run:   Run,
}

var branchName string

var baseURL string = "https://api.github.com/repos/Minestom/Minestom/commits/"

func Execute() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	root.PersistentFlags().StringVarP(&branchName, "branch", "b", "master", "Example: -b 1_20_5 would pick the 1_20_5 branch")
}

func Run(cobra *cobra.Command, args []string) {
	color.Cyan("Getting latest Minestom version for branch " + branchName)

	GetCommit(branchName)
}

func GetCommit(id string) {
	resp, err := http.Get(baseURL + id)
	if err != nil {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Printf("ERROR: Error making request %v", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		HandleOK(resp)
		return
	case http.StatusNotFound:
		Handle404(resp)
		return
	case http.StatusForbidden:
		c := color.New(color.FgRed).Add(color.Bold)
		var apiResp util.RateLimit
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			c.Printf("ERROR: Error reading body %v", err)
			os.Exit(1)
		}
		c.Printf("ERROR: Forbidden %v", apiResp.Message)
		return
	}

}

func SuccessfulCommit(data util.GHSuccessResponse) bool {

	resp, err := http.Get(baseURL + data.Sha + "/check-runs")
	if err != nil {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Printf("ERROR: Error making request %v", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var apiResp util.GHCheckRunsSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			c := color.New(color.FgRed).Add(color.Bold)
			c.Printf("ERROR: Error reading body %v", err)
			os.Exit(1)
		}

		for i := 0; i < len(apiResp.CheckRuns); i++ {
			var run = apiResp.CheckRuns[i]

			if run.Conclusion != "success" {
				return false
			}
		}

		return true
	case http.StatusNotFound:
		return false
	}

	return false

}

func HandleOK(resp *http.Response) {
	var apiResp util.GHSuccessResponse
	err := json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Printf("ERROR: Error reading body %v", err)
		os.Exit(1)
	}

	var sha = apiResp.Sha

	c := color.New(color.FgYellow)
	c.Println("Checking if commit " + sha + " is successful")
	var successfulCommit = SuccessfulCommit(apiResp)

	if !successfulCommit {
		c = color.New(color.FgRed)
		c.Println("Commit " + sha + " is not succesful")
		GetCommit(apiResp.Parents[0].Sha)
		return
	}

	c = color.New(color.FgGreen)
	c.Println(
		sha[0:10]+" ("+sha+")",
		"-",
		apiResp.Commit.Message,
		"-",
		apiResp.Commit.Author.Name,
	)

}

func Handle404(resp *http.Response) {
	var apiResp util.GHErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&apiResp)
	c := color.New(color.FgRed).Add(color.Bold)
	if err != nil {
		c.Printf("ERROR: Error reading body %v", err)
		os.Exit(1)
	}

	c.Printf("ERROR: %v", apiResp.Message)
}
