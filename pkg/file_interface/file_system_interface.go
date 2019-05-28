package fileinterface

import (
	"io/ioutil"
	"os"

	yamlresources "github.com/redhat-cop/applier-cli/pkg/yaml_resources"
	yaml "gopkg.in/yaml.v2"
)

// FileSystemInterface is an implementation of FileInterface which interacts with the filesystem on the local machine.
type FileSystemInterface struct{}

// Mkdir creates a directory on the local filesystem.
func (fileInterface *FileSystemInterface) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// ReadClusterContents parses and unmarshals the group_vars/all.yml file, which contains openshift_cluster_content.
func (fileInterface *FileSystemInterface) ReadClusterContents() (yamlresources.ClusterContentList, error) {
	var clusterContents yamlresources.ClusterContentList
	fileContents, err := ioutil.ReadFile("inventory/group_vars/all.yml")
	if err != nil {
		return clusterContents, err
	}
	err = yaml.Unmarshal(fileContents, &clusterContents)
	return clusterContents, err
}

// ReadFile simply reads a file from the filesystem and returns the content as a byte array
func (fileInterface *FileSystemInterface) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// ReadResource reads and unmarshals a yaml file from the filesystem.
func (fileInterface *FileSystemInterface) ReadResource(fileName string) (map[string]interface{}, error) {
	var resource map[string]interface{}

	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return resource, err
	}
	err = yaml.Unmarshal(fileContents, &resource)

	return resource, err
}

// TouchParamsFile creates the params file with a template if it does not exist, and writes an example to it.
func (fileInterface *FileSystemInterface) TouchParamsFile(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(fileName, []byte("# Use this parameter file as shown:\n# PARAMETER=value"), 0766)
	}
	return err
}

// WriteClusterContents writes the openshift_cluster_content to the group_vars/all.yml file.
func (fileInterface *FileSystemInterface) WriteClusterContents(contents yamlresources.ClusterContentList) error {
	outByte := []byte{}
	outByte, _ = yaml.Marshal(contents)

	err := ioutil.WriteFile("inventory/group_vars/all.yml", outByte, 0766)
	return err
}

// WriteFile writes the data given to the path given on the local filesystem.
func (fileInterface *FileSystemInterface) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
