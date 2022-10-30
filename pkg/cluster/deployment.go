package cluster

import (
	"context"
	"fmt"
)

func newDeploymentTask() clusterTask {
	return deploymentTask{}
}

type deploymentTask struct{}

func (deploymentTask) Run(ctx context.Context, config ClusterReconcilerConfig) error {
	return fmt.Errorf("not implemented")
}
