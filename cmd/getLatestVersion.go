package cmd

import (
	"fmt"
	"log"

	githubapi "github.com/redhat-cop/applier-cli/pkg/github_api"
	"github.com/spf13/cobra"
)

// getLatestVersionCmd represents the get-latest-version command
var getLatestVersionCmd = &cobra.Command{
	Use:   "get-latest-version",
	Short: "Display the latest version of OpenShift-Applier",
	Long: `Uses the GitHub.com API to determine the latest released version
of OpenShift-Applier. For information purposes only, and not a prerequisite
for any other command (init performs this action automatically).`,
	Run: func(cmd *cobra.Command, args []string) {
		releaseAPI := githubapi.GitHubReleaseAPI{}
		getLatestVer(&releaseAPI)
	},
}

func getLatestVer(api githubapi.ReleaseAPI) {
	release, err := api.GetLatestVersionInfo()
	if err != nil {
		log.Fatal("Could not fetch the latest version of OpenShift-Applier.")
	}
	fmt.Println(release.TagName)
}

func init() {
	rootCmd.AddCommand(getLatestVersionCmd)
}
