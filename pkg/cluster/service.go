package cluster

import (
	"context"
	"fmt"
)

func newServiceTask() clusterTask {
	return serviceTask{}
}

type serviceTask struct{}

func (serviceTask) Run(ctx context.Context, config ClusterReconcilerConfig) error {
	return fmt.Errorf("not implemented")
}
