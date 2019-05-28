package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the current OpenShift-Applier inventory.",
	Long: `Runs the current OpenShift-Applier inventory. Default is 
to run using local Ansible, but can also run OpenShift-Applier 
in a Docker container.`,
	Run: func(cmd *cobra.Command, args []string) {
		docker, _ := cmd.Flags().GetBool("docker")

		if docker {
			var SELinuxModifier string
			if checkSELinux() {
				SELinuxModifier = ":z"
			} else {
				SELinuxModifier = ""
			}
			user, err := user.Current()
			if err != nil {
				log.Fatal("Could not determine current user ID")
			}
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatal("Could not determine current working directory")
			}
			docker := exec.Command(
				"docker",
				"run",
				"-u", user.Uid,
				"-v", fmt.Sprintf("%s/.kube/config:/openshift-applier/.kube/config%s", user.HomeDir, SELinuxModifier),
				"-v", fmt.Sprintf("%s/:/tmp/applier%s", pwd, SELinuxModifier),
				"-e", "INVENTORY_PATH=/tmp/applier/inventory",
				"-t", "redhatcop/openshift-applier",
			)
			docker.Dir = pwd
			docker.Stdout = os.Stdout
			docker.Stderr = os.Stderr
			if err := docker.Run(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Finished executing Ansible playbook in Docker container.")
		} else {
			ansible := exec.Command("ansible-playbook", "./apply.yml", "-i", "inventory/")
			ansible.Dir, _ = os.Getwd()
			ansible.Stdout = os.Stdout
			ansible.Stderr = os.Stderr
			ansible.Env = append(os.Environ(),
				"OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES", // Attempt to fix thread safety issue when running natively on macOS. Ref: https://github.com/ansible/ansible/issues/32499
			)
			if err := ansible.Run(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Finished executing Ansible playbook.")
		}
	},
}

func checkSELinux() bool {
	var out bytes.Buffer
	cmd := exec.Command("getenforce")
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil && strings.Trim(out.String(), "\n") == "Enforcing" {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("docker", "d", false, "Run using the OpenShift-Applier Docker image instead of local Ansible.")
}
