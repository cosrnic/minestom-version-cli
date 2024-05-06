package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cosrnic/minestom-version-cli/util"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "minestom-version",
	Short: "Minestom-version CLI in Go",
	Long:  "A CLI for getting the latest version of Minestom [or a branch of minestom]",
	Run:   Run,
}

var branchName string

var baseURL string = "https://api.github.com/repos/Minestom/Minestom/branches/"

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
	fmt.Println("Getting latest Minestom version for branch " + branchName)

	resp, err := http.Get(baseURL + branchName)
	if err != nil {
		fmt.Printf("ERROR: Error making request %v", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		HandleOK(resp)
	case http.StatusNotFound:
		Handle404(resp)
	}

}

func HandleOK(resp *http.Response) {
	var apiResp util.GHSuccessResponse
	err := json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Printf("ERROR: Error reading body %v", err)
		os.Exit(1)
	}

	fmt.Println(
		apiResp.Commit.Sha[0:10]+" ("+apiResp.Commit.Sha+")",
		"-",
		apiResp.Commit.Commit.Message,
		"-",
		apiResp.Commit.Commit.Author.Name,
	)

}

func Handle404(resp *http.Response) {
	var apiResp util.GHErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Printf("ERROR: Error reading body %v", err)
		os.Exit(1)
	}

	fmt.Printf("ERROR: %v", apiResp.Message)
}
