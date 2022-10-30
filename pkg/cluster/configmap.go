package cluster

import (
	"context"
	"fmt"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func newConfigmapTask() clusterTask {
	return configmapTask{}
}

type configmapTask struct{}

func (t configmapTask) Run(ctx context.Context, config ClusterReconcilerConfig) error {
	var (
		err      error
		desired  corev1.ConfigMap
		existing *corev1.ConfigMap = &corev1.ConfigMap{}
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
	return err
}

func (configmapTask) genDesired(ins ksqldbv1alpha1.KsqldbCluster) corev1.ConfigMap {
	//TODO: Implement
	return corev1.ConfigMap{}
}
