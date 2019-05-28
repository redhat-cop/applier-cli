package cmd

import (
	"testing"

	clusterinterface "github.com/redhat-cop/applier-cli/pkg/cluster_interface"
	fileinterface "github.com/redhat-cop/applier-cli/pkg/file_interface"
	githubapi "github.com/redhat-cop/applier-cli/pkg/github_api"
	"github.com/stretchr/testify/assert"
)

func TestAddFromClusterMakeTemplate(t *testing.T) {

	// Test adding a resource from a cluster, turning it into a template
	testFlags := runFlags{
		fromCluster:  true,
		fromFile:     false,
		makeTemplate: true,
		edit:         false,
	}
	clusterInterface := clusterinterface.MockClusterInterface{}
	fileInterface := fileinterface.MockFileInterface{}
	releaseAPI := githubapi.MockReleaseAPI{}
	initRun(&fileInterface, &releaseAPI)
	add(testFlags, []string{"test_pod"}, &clusterInterface, &fileInterface)

	_, err := fileInterface.ReadFile("templates/nginx-pod.yml")
	assert.NoError(t, err)
	_, err = fileInterface.ReadFile("params/nginx-pod")
	assert.NoError(t, err)

}

func TestAddFromClusterNotTemplate(t *testing.T) {

	// Test adding a resource from a cluster, not turning it into a template
	testFlags := runFlags{
		fromCluster:  true,
		fromFile:     false,
		makeTemplate: false,
		edit:         false,
	}
	clusterInterface := clusterinterface.MockClusterInterface{}
	fileInterface := fileinterface.MockFileInterface{}
	releaseAPI := githubapi.MockReleaseAPI{}
	initRun(&fileInterface, &releaseAPI)
	add(testFlags, []string{"test_pod"}, &clusterInterface, &fileInterface)

	_, err := fileInterface.ReadFile("files/nginx-pod.yml")
	assert.NoError(t, err)

}

func TestAddFromFileMakeTemplate(t *testing.T) {

	// Test adding a resource from a cluster, not turning it into a template
	testFlags := runFlags{
		fromCluster:  false,
		fromFile:     true,
		makeTemplate: true,
		edit:         false,
	}
	clusterInterface := clusterinterface.MockClusterInterface{}
	fileInterface := fileinterface.MockFileInterface{}
	releaseAPI := githubapi.MockReleaseAPI{}
	initRun(&fileInterface, &releaseAPI)
	// Initialize test file in mock FS
	fileInterface.WriteFile("some_dir/test_resource.yml", []byte(`apiVersion: v1
kind: Pod
metadata:
  name: nginx
  annotations:
    test: test
spec:
  containers:
  - name: nginx
    image: nginx
`), 766)
	add(testFlags, []string{"some_dir/test_resource.yml"}, &clusterInterface, &fileInterface)

	_, err := fileInterface.ReadResource("templates/nginx-pod.yml")
	assert.NoError(t, err)
	_, err = fileInterface.ReadFile("params/nginx-pod")
	assert.NoError(t, err)

}

func TestAddFromFileNotTemplate(t *testing.T) {

	// Test adding a resource from a cluster, not turning it into a template
	testFlags := runFlags{
		fromCluster:  false,
		fromFile:     true,
		makeTemplate: false,
		edit:         false,
	}
	clusterInterface := clusterinterface.MockClusterInterface{}
	fileInterface := fileinterface.MockFileInterface{}
	releaseAPI := githubapi.MockReleaseAPI{}
	initRun(&fileInterface, &releaseAPI)
	// Initialize test file in mock FS
	fileInterface.WriteFile("some_dir/test_resource.yml", []byte(`apiVersion: v1
kind: Pod
metadata:
  name: nginx
  annotations:
    test: test
spec:
  containers:
  - name: nginx
    image: nginx
`), 766)
	add(testFlags, []string{"some_dir/test_resource.yml"}, &clusterInterface, &fileInterface)

	_, err := fileInterface.ReadResource("files/nginx-pod.yml")
	assert.NoError(t, err)

}
