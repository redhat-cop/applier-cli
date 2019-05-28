package githubapi

// ReleaseAPI represents the interface for determining the latest version of OpenShift-Applier
type ReleaseAPI interface {
	GetLatestVersionInfo() (GitHubRelease, error)
}
