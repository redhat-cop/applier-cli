package fileinterface

import (
	"errors"
	"os"

	yamlresources "github.com/redhat-cop/applier-cli/pkg/yaml_resources"
	yaml "gopkg.in/yaml.v2"
)

// MockFileInterface is an implementation of FileInterface which works entirely in-memory for testing.
type MockFileInterface struct {
	ClusterContents []byte
	Files           map[string][]byte
}

// Mkdir adds an element signalling a directory on the virtual filesystem
func (fileInterface *MockFileInterface) Mkdir(name string, perm os.FileMode) error {
	return fileInterface.WriteFile(name, []byte("[directory]"), perm)
}

// ReadClusterContents parses and unmarshals the virtual file containing openshift_cluster_content.
func (fileInterface *MockFileInterface) ReadClusterContents() (yamlresources.ClusterContentList, error) {
	var ClusterContents yamlresources.ClusterContentList
	err := yaml.Unmarshal(fileInterface.ClusterContents, &ClusterContents)
	return ClusterContents, err
}

// ReadFile simply reads a file from the virtual filesystem and returns the content as a byte array
func (fileInterface *MockFileInterface) ReadFile(filename string) ([]byte, error) {
	if contents, ok := fileInterface.Files[filename]; ok {
		return contents, nil
	}
	return []byte{}, errors.New("could not find file")
}

// ReadResource reads and unmarshals a yaml file from the virtual filesystem.
func (fileInterface *MockFileInterface) ReadResource(filename string) (map[string]interface{}, error) {
	var resource map[string]interface{}

	// fileContents := fileInterface.Files[filename]
	if contents, ok := fileInterface.Files[filename]; ok {
		err := yaml.Unmarshal(contents, &resource)
		return resource, err
	}
	return map[string]interface{}{}, errors.New("could not find file")
}

// TouchParamsFile creates the params file for a template on the virtual filesystem, and writes an example to it.
func (fileInterface *MockFileInterface) TouchParamsFile(filename string) error {
	err := fileInterface.WriteFile(filename, []byte("# Use this parameter file as shown:\n# PARAMETER=value"), 0666)
	return err
}

// WriteClusterContents writes the openshift_cluster_content to the virtual filesystem.
func (fileInterface *MockFileInterface) WriteClusterContents(contents yamlresources.ClusterContentList) error {
	outByte := []byte{}
	outByte, err := yaml.Marshal(contents)
	fileInterface.ClusterContents = outByte
	return err
}

// WriteFile writes the data given to the path given on the virtual filesystem
func (fileInterface *MockFileInterface) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if fileInterface.Files == nil {
		fileInterface.Files = map[string][]byte{}
	}
	fileInterface.Files[filename] = data
	return nil
}
