package fileinterface

import (
	"os"

	yamlresources "github.com/redhat-cop/applier-cli/pkg/yaml_resources"
)

// FileInterface defines methods for interacting with files on the local filesystem
type FileInterface interface {
	Mkdir(name string, perm os.FileMode) error
	ReadClusterContents() (yamlresources.ClusterContentList, error)
	ReadFile(filename string) ([]byte, error)
	ReadResource(fileName string) (map[string]interface{}, error)
	TouchParamsFile(fileName string) error
	WriteClusterContents(yamlresources.ClusterContentList) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
}
