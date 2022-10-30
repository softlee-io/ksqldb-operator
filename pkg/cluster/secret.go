package cluster

import (
	"context"
	"fmt"
)

func newSecretTask() clusterTask {
	return secretTask{}
}

type secretTask struct{}

func (secretTask) Run(ctx context.Context, config ClusterReconcilerConfig) error {
	return fmt.Errorf("not implemented")
}
