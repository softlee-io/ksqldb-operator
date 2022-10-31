package cluster

import (
	"context"
	"fmt"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func newDeploymentTask() clusterTask {
	return deploymentTask{}
}

type deploymentTask struct{}

func (t deploymentTask) Run(ctx context.Context, config ClusterReconcilerConfig) error {
	var (
		err      error
		desired  appsv1.Deployment
		existing *appsv1.Deployment = &appsv1.Deployment{}
	)

	desired = t.genDesired(config.Instance)
	err = controllerutil.SetControllerReference(&config.Instance, &desired, config.Scheme)
	if err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	nsName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
	err = config.Client.Get(ctx, nsName, existing)

	//TODO:
	//- Route through ErrorHnadling (doesExist)
	//- IsNotFound error
	//	- createDesired
	//	- compareWithDesired
	//- Other Errors
	//	- returen directly
	//- IfFound
	//	- patchResource
	return fmt.Errorf("not implemented")
}

func (t deploymentTask) genDesired(ins ksqldbv1alpha1.KsqldbCluster) appsv1.Deployment {
	//TODO: Implement
	return appsv1.Deployment{}
}
