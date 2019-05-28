package githubapi

import (
	"encoding/json"
	"net/http"
	"time"
)

// GitHubReleaseAPI represents the object which actually connects to GitHub to determine the latest release version
type GitHubReleaseAPI struct {
}

// GetLatestVersionInfo calls the GitHub API and returns a GitHubRelease object representing the latest release of OpenShift-Applier
func (releaseAPI *GitHubReleaseAPI) GetLatestVersionInfo() (GitHubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	response := &GitHubRelease{}
	r, err := client.Get("https://api.github.com/repos/redhat-cop/openshift-applier/releases/latest")
	if err != nil {
		return *response, err
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(response)
	return *response, err
}
