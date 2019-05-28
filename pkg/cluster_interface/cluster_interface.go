package clusterinterface

// ClusterInterface defines methods for interacting with a k8s cluster
type ClusterInterface interface {
	GetResource(resourceName string) (map[string]interface{}, error)
}
