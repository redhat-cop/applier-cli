package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	fileinterface "github.com/redhat-cop/applier-cli/pkg/file_interface"
	githubapi "github.com/redhat-cop/applier-cli/pkg/github_api"
)

func TestInit(t *testing.T) {

	releaseAPI := githubapi.MockReleaseAPI{}
	fileInterface := fileinterface.MockFileInterface{}
	initRun(&fileInterface, &releaseAPI)

	assertExistsOnFileSystem := []string{
		"inventory",
		"inventory/host_vars",
		"inventory/host_vars/localhost.yml",
		"inventory/group_vars",
		"inventory/group_vars/all.yml",
		"inventory/hosts",
		"templates",
		"params",
		"files",
		"requirements.yml",
		"apply.yml",
	}

	for _, val := range assertExistsOnFileSystem {
		_, err := fileInterface.ReadFile(val)
		assert.NoError(t, err)
	}

}
