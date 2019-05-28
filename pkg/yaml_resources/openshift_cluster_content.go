package yamlresources

// ClusterContentList represents openshift_cluster_content within an OpenShift-Applier inventory
type ClusterContentList struct {
	OpenShiftClusterContent []ClusterContentObject `yaml:"openshift_cluster_content"`
}

// ClusterContentObject represents a single object within ClusterContentList
type ClusterContentObject struct {
	Object  string           `yaml:"object"`
	Content []ClusterContent `yaml:"content"`
}

// ClusterContent represents the actual content of a ClusterContentObject
type ClusterContent struct {
	Name           string `yaml:"name"`
	File           string `yaml:"file,omitempty"`
	Template       string `yaml:"template,omitempty"`
	Params         string `yaml:"params,omitempty"`
	ParamsFromVars string `yaml:"params_from_vars,omitempty"`
	Action         string `yaml:"action,omitempty"`
}
