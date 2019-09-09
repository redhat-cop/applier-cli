package cmd

import (
	"fmt"
	"log"
	"strings"

	fileinterface "github.com/redhat-cop/applier-cli/pkg/file_interface"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// getCurrentVersionCmd represents the get-current-version command
var getCurrentVersionCmd = &cobra.Command{
	Use:   "get-current-version",
	Short: "Display the current version of the OpenShift-Applier",
	Long:  `Reads the current version of the applier defined in requirements.yml`,

	Run: func(cmd *cobra.Command, args []string) {
		var fileInterface *fileinterface.FileSystemInterface
		getCurrentVer(fileInterface)
	},
}

func getCurrentVer(fileInterface fileinterface.FileInterface) {

	reqFile, err := fileInterface.ReadFile("./requirements.yml")
	if err != nil {
		return
	}

	var reqs []map[string]interface{}
	err = yaml.Unmarshal(reqFile, &reqs)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range reqs {
		if strings.ToLower(e["name"].(string)) == "openshift-applier" {
			fmt.Println(e["version"])
			break
		} else {
			log.Fatal("Cannot find openshift-applier in your requirements file.")
		}
	}
}

func init() {
	rootCmd.AddCommand(getCurrentVersionCmd)
}
