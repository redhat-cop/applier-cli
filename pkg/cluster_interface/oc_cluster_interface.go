package clusterinterface

import (
	"bytes"
	"log"
	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

// OCClusterInterface is an implementation of ClusterInterface which uses the local installation of `oc` for interacting with the cluster
type OCClusterInterface struct {
}

// GetResource uses the `oc` CLI to export a resource from the cluster and convert it into a map[string]interface{}
func (cluster *OCClusterInterface) GetResource(resourceName string) (map[string]interface{}, error) {
	var out bytes.Buffer
	var resource map[string]interface{}

	cmd := exec.Command("oc", "get", resourceName, "--export", "-o", "yaml")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("Could not export the desired resource from the cluster. Are you sure that you are logged in and it exists?")
	}
	err = yaml.Unmarshal(out.Bytes(), &resource)
	if err != nil {
		log.Fatal("Unable to interpret the resource.")
	}

	return resource, err
}
