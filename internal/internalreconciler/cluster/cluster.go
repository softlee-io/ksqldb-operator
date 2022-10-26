package cluster

import (
	"context"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/config"
)

type ClusterReconciler struct {
	config.BaseParam
	Instance ksqldbv1alpha1.KsqldbCluster
}

func (r ClusterReconciler) Start(ctx context.Context) error {
	//r.Instance.
	return nil
}

func (r ClusterReconciler) deployment(ctx context.Context) error {
	//desired :=

	return nil
}
