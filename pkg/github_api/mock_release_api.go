package githubapi

// MockReleaseAPI represents the object which actually connects to GitHub to determine the latest release version
type MockReleaseAPI struct {
	TagName string
}

// GetLatestVersionInfo calls the GitHub API and returns a GitHubRelease object representing the latest release of OpenShift-Applier
func (releaseAPI *MockReleaseAPI) GetLatestVersionInfo() (GitHubRelease, error) {
	if releaseAPI.TagName == "" {
		releaseAPI.TagName = "1.0.0"
	}
	response := &GitHubRelease{
		TagName: releaseAPI.TagName,
	}
	return *response, nil
}
