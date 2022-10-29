package query

import (
	"context"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/config"
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
// - ksqldb api calls towards selected ksqldb cluste
func (r queryReconciler) Start(ctx context.Context) error {
	return nil
}
