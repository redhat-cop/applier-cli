package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	clusterinterface "github.com/redhat-cop/applier-cli/pkg/cluster_interface"
	fileinterface "github.com/redhat-cop/applier-cli/pkg/file_interface"
	yamlresources "github.com/redhat-cop/applier-cli/pkg/yaml_resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type runFlags struct {
	fromCluster  bool
	fromFile     bool
	makeTemplate bool
	edit         bool
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new resources to the current inventory",
	Long: `Add new resources to the current inventory. Resources that
are templates are added to the templates directory, while all
others are added to the files directory. Non-template resources
can also be converted into templates. Generated resources can immediately
be opened in your default editor for further tuning.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get all flags
		fromCluster, _ := cmd.Flags().GetBool("from-cluster")
		fromFile, _ := cmd.Flags().GetBool("from-file")
		makeTemplate, _ := cmd.Flags().GetBool("make-template")
		edit, _ := cmd.Flags().GetBool("edit")

		flags := runFlags{
			fromCluster:  fromCluster,
			fromFile:     fromFile,
			makeTemplate: makeTemplate,
			edit:         edit,
		}

		var clusterInterface *clusterinterface.OCClusterInterface
		var fileInterface *fileinterface.FileSystemInterface

		add(flags, args, clusterInterface, fileInterface)
	},
}

func add(flags runFlags, args []string, clusterInterface clusterinterface.ClusterInterface, fileInterface fileinterface.FileInterface) {
	var resource map[string]interface{}
	var err error

	if flags.fromCluster {
		resource, err = clusterInterface.GetResource(args[0])
	} else if flags.fromFile {
		resource, err = fileInterface.ReadResource(args[0])
	} else {
		log.Fatal("Unclear where to get the resource. Please use --from-cluster (-c) or --from-file (-f).")
	}
	if err != nil {
		log.Fatal("Unable to get the resource requested.", err)
	}

	resourceKind := resource["kind"].(string)
	name := resource["metadata"].(map[interface{}]interface{})["name"].(string)
	var fileCreated string

	filterUndesireableKeys(resource)

	if flags.makeTemplate || resourceKind == "Template" {
		fileCreated = writeTemplateToInventory(resource, resourceKind, name, fileInterface)
	} else {
		fileCreated = writeFileToInventory(resource, resourceKind, name, fileInterface)
	}

	if flags.edit {
		cmd := exec.Command(viper.GetString("editor"), fileCreated)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func filterUndesireableKeys(resource map[string]interface{}) {
	// anything more than one level in needs to have a check to see if the outer level exists.
	if _, hasMetadata := resource["metadata"]; hasMetadata {
		delete(resource["metadata"].(map[interface{}]interface{}), "selfLink")
		delete(resource["metadata"].(map[interface{}]interface{}), "generation")
		delete(resource["metadata"].(map[interface{}]interface{}), "creationTimestamp")
		delete(resource["metadata"].(map[interface{}]interface{}), "resourceVersion")
		if _, hasAnnotations := resource["metadata"].(map[interface{}]interface{})["annotations"]; hasAnnotations {
			delete(resource["metadata"].(map[interface{}]interface{})["annotations"].(map[interface{}]interface{}), "kubectl.kubernetes.io/last-applied-configuration")
			delete(resource["metadata"].(map[interface{}]interface{})["annotations"].(map[interface{}]interface{}), "deployment.kubernetes.io/revision")
		}
	}
	delete(resource, "status")
}

func writeTemplateToInventory(resource map[string]interface{}, kind string, name string, fileInterface fileinterface.FileInterface) string {

	outByte := []byte{}
	lowerCaseKind := strings.ToLower(kind)

	if kind != "Template" {
		template := yamlresources.Template{
			APIVersion: "v1",
			Kind:       "Template",
			Metadata: yamlresources.TemplateMetadata{
				Name: name,
			},
			Objects: []map[string]interface{}{
				resource,
			},
			Parameters: nil,
		}
		outByte, _ = yaml.Marshal(template)
	} else {
		outByte, _ = yaml.Marshal(resource)
	}

	outFile := fmt.Sprintf("templates/%s-%s.yml", name, lowerCaseKind)
	paramsFile := fmt.Sprintf("params/%s-%s", name, lowerCaseKind)
	err := fileInterface.WriteFile(outFile, outByte, 0766)

	fileInterface.TouchParamsFile(paramsFile)

	clusterContents, err := fileInterface.ReadClusterContents()
	if err != nil {
		log.Fatal("Unable to read the inventory in the current directory. Check the path or run \"applier-cli init\".")
	}
	clusterContents.OpenShiftClusterContent = append(clusterContents.OpenShiftClusterContent, yamlresources.ClusterContentObject{
		Object: name,
		Content: []yamlresources.ClusterContent{{
			Name:     name,
			Template: fmt.Sprintf("{{ inventory_dir }}/../%s", outFile),
			Params:   fmt.Sprintf("{{ inventory_dir }}/../%s", paramsFile),
		}},
	})
	err = fileInterface.WriteClusterContents(clusterContents)

	if err != nil {
		log.Fatal("Unable to add the template to the current inventory.")
	} else {
		fmt.Println("Template added to the current inventory.")
	}

	return outFile

}

func writeFileToInventory(resource map[string]interface{}, kind string, name string, fileInterface fileinterface.FileInterface) string {

	outByte, _ := yaml.Marshal(resource)
	lowerCaseKind := strings.ToLower(kind)

	outFile := fmt.Sprintf("files/%s-%s.yml", name, lowerCaseKind)
	err := fileInterface.WriteFile(outFile, outByte, 0766)
	if err != nil {
		log.Fatal("Unable to write the file to the current inventory.")
	}

	clusterContents, err := fileInterface.ReadClusterContents()
	if err != nil {
		log.Fatal("Unable to read the inventory in the current directory. Check the path or run \"applier-cli init\".")
	}

	clusterContents.OpenShiftClusterContent = append(clusterContents.OpenShiftClusterContent, yamlresources.ClusterContentObject{
		Object: name,
		Content: []yamlresources.ClusterContent{{
			Name: name,
			File: fmt.Sprintf("{{ inventory_dir }}/../%s", outFile),
		}},
	})
	fileInterface.WriteClusterContents(clusterContents)

	if err != nil {
		log.Fatal(err)
		fmt.Println("Could not add file to the current inventory.")
	} else {
		fmt.Println("File added to the current inventory.")
	}

	return outFile

}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("from-cluster", "c", false, "Use an existing cluster as a source for this resource")
	addCmd.Flags().BoolP("from-file", "f", false, "Use a yaml file as a source for this resource")
	addCmd.Flags().BoolP("make-template", "t", false, "Convert the resource into a template (if it isn't already)")
	addCmd.Flags().BoolP("edit", "e", false, "Immediately open the file in your default editor once created")
}
