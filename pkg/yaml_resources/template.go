package yamlresources

// Template represents an OpenShift template resource
type Template struct {
	APIVersion string                   `yaml:"apiVersion"`
	Kind       string                   `yaml:"kind"`
	Metadata   TemplateMetadata         `yaml:"metadata"`
	Objects    []map[string]interface{} `yaml:"objects"`
	Parameters []TemplateParameter      `yaml:"parameters"`
}

// TemplateMetadata represents the metadata within an OpenShift template resource
type TemplateMetadata struct {
	Name string `yaml:"name"`
}

// TemplateParameter represents a single parameter within an OpenShift template resource
type TemplateParameter struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
}
