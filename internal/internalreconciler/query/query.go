package query

import (
	"context"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/config"
)

type QueryReconciler struct {
	config.BaseParam
	Instance ksqldbv1alpha1.KsqldbQuery
}

func (r QueryReconciler) Start(ctx context.Context) error {
	return nil
}
