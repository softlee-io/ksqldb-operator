package query

import (
	"context"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
)

type QueryReconcilerConfig struct {
	config.BaseParam
	Instance ksqldbv1alpha1.KsqldbQuery
}

func (c QueryReconcilerConfig) Validate() error {
	//TODO: Impl validation logic
	return nil
}

type queryReconciler struct {
	config QueryReconcilerConfig
}

func NewQueryReconciler(config QueryReconcilerConfig) queryReconciler {
	return queryReconciler{
		config: config,
	}
}

// No provision task included:
// - ksqldb api calls towards selected ksqldb cluster (Pod DNS)
// - desiredState = REST Get against one of cluster pods
//   - Save permanent push-based query as configmap
//   - Create: If difference is detected, send update request ()
//   - Update
//
// -
func (r queryReconciler) Start(ctx context.Context) error {
	return nil
}
