package clusterinterface

// MockClusterInterface is an implementation of ClusterInterface which only returns static data
type MockClusterInterface struct {
}

// GetResource uses the `oc` CLI to export a resource from the cluster and convert it into a map[string]interface{}
func (cluster *MockClusterInterface) GetResource(resourceName string) (map[string]interface{}, error) {

	var dummyData map[string]interface{}
	switch resourceName {
	case "test_pod":
		dummyData = map[string]interface{}{
			"kind": "pod",
			"metadata": map[interface{}]interface{}{
				"name":        "nginx",
				"annotations": map[interface{}]interface{}{},
			},
		}
	}
	return dummyData, nil

}
