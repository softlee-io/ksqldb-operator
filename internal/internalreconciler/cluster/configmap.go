package cluster

import "fmt"

func newConfigmapTask() clusterTask {
	return configmapTask{}
}

type configmapTask struct{}

func (configmapTask) run(clusterRecon ClusterReconcilerConfig) error {
	return fmt.Errorf("not implemented")
}
