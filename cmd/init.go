package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	fileinterface "github.com/redhat-cop/applier-cli/pkg/file_interface"
	githubapi "github.com/redhat-cop/applier-cli/pkg/github_api"
	yamlresources "github.com/redhat-cop/applier-cli/pkg/yaml_resources"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an empty OpenShift-Applier inventory",
	Long: `Scaffolds an empty OpenShift-Applier inventory, including:

An inventory directory, with host_vars and group_vars
A templates directory
A params directory
A files directory
An Ansible Galaxy requirements file with the latest release of OpenShift-Applier on GitHub
	
In addition, the Ansible Galaxy requirements are installed if the ansible-galaxy bin is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileInterface *fileinterface.FileSystemInterface
		var releaseAPI *githubapi.GitHubReleaseAPI
		initRun(fileInterface, releaseAPI)
	},
}

func initRun(fileInterface fileinterface.FileInterface, releaseAPI githubapi.ReleaseAPI) {
	makeAllDirectories(fileInterface)
	latestReleasedVersion := getLatestApplierReleaseTag(releaseAPI)
	writeConfigs(fileInterface)
	writeGalaxyRequirements(latestReleasedVersion, fileInterface)
	installGalaxyRequirements()
}

func makeAllDirectories(fileInterface fileinterface.FileInterface) {
	fileInterface.Mkdir("inventory", 0766)
	fileInterface.Mkdir("inventory/host_vars", 0766)
	fileInterface.Mkdir("inventory/group_vars", 0766)
	fileInterface.Mkdir("templates", 0766)
	fileInterface.Mkdir("params", 0766)
	fileInterface.Mkdir("files", 0766)
}

func writeConfigs(fileInterface fileinterface.FileInterface) {
	hosts := []byte("[seed-hosts]\nlocalhost")
	err := fileInterface.WriteFile("inventory/hosts", hosts, 0766)
	if err != nil {
		log.Fatal("Could not write inventory/hosts")
	}

	hostVars := []byte("ansible_connection: local")
	err = fileInterface.WriteFile("inventory/host_vars/localhost.yml", hostVars, 0766)
	if err != nil {
		log.Fatal("Could not write inventory/host_vars/localhost.yml")
	}

	ansiblePlaybook := []byte(`- name: Deploy OpenShift-Applier Inventory
  hosts: seed-hosts[0]
  tasks:
  - include_role:
      name: roles/openshift-applier/roles/openshift-applier
    tags:
    - openshift-applier`)
	err = fileInterface.WriteFile("apply.yml", ansiblePlaybook, 0766)
	if err != nil {
		log.Fatal("Could not write apply.yml")
	}

	groupVars := &yamlresources.ClusterContentList{}
	yamlGroupVars, _ := yaml.Marshal(groupVars)
	err = fileInterface.WriteFile("inventory/group_vars/all.yml", yamlGroupVars, 0766)
	if err != nil {
		log.Fatal("Could not write inventory/group_vars/all.yml")
	}
}

func getLatestApplierReleaseTag(releaseAPI githubapi.ReleaseAPI) string {
	currentApplierVersion, err := releaseAPI.GetLatestVersionInfo()
	if err != nil {
		log.Fatal("Could not determine the latest release of OpenShift-Applier.")
	}
	return currentApplierVersion.TagName
}

func writeGalaxyRequirements(version string, fileInterface fileinterface.FileInterface) {
	requirements := &yamlresources.Requirements{{
		Src:     "https://github.com/redhat-cop/openshift-applier",
		SCM:     "git",
		Version: version,
		Name:    "openshift-applier",
	}}
	yamlRequirements, err := yaml.Marshal(requirements)
	if err != nil {
		log.Fatal("Could not generate requirements file.")
	}
	err = fileInterface.WriteFile("requirements.yml", yamlRequirements, 0766)
	if err != nil {
		log.Fatal("Could not write requirements file.")
	}
}

func installGalaxyRequirements() {
	galaxy := exec.Command("ansible-galaxy", "install", "-r", "requirements.yml", "--roles-path=roles", "-f")
	galaxy.Dir, _ = os.Getwd()
	_, err := galaxy.Output()
	fmt.Println("Initialized empty openshift-applier inventory")
	if err != nil {
		fmt.Println("Could not invoke ansible-galaxy to install requirements. Please check your installation and install ansible-galaxy requirements manually.")
	} else {
		fmt.Println("Successfully installed ansible-galaxy requirements")
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
