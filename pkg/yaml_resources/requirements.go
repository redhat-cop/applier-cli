package yamlresources

// Requirements represents the data in requirements.yaml, which are Ansible Galaxy requirements.
type Requirements []struct {
	Src     string `yaml:"src"`
	SCM     string `yaml:"scm"`
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
}
