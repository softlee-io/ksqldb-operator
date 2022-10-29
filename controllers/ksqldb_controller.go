/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/cluster"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/config"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/query"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconObjectType string

const (
	KSQLDB_CLUSTER ReconObjectType = "KsqldbCluster"
	KSQLDB_QUERY   ReconObjectType = "KsqldbQuery"
)

// KsqldbReconciler reconciles a Ksqldb objects
type KsqldbReconciler struct {
	base       config.BaseParam
	reconOrder []ReconObjectType
}

func NewReconciler(param config.BaseParam) *KsqldbReconciler {
	objOrders := []ReconObjectType{KSQLDB_CLUSTER, KSQLDB_QUERY}
	return &KsqldbReconciler{
		base:       param,
		reconOrder: objOrders,
	}
}

//+kubebuilder:rbac:groups=ksqldb.softlee.io,resources=ksqldbclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ksqldb.softlee.io,resources=ksqldbclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ksqldb.softlee.io,resources=ksqldbclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ksqldb objects against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KsqldbReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	arrErr := []error{}

	var curErr error
	for i := 0; i < len(r.reconOrder); i++ {
		if i != 0 && curErr == nil {
			switch r.reconOrder[i] {
			case KSQLDB_CLUSTER:
				var clusterInstance v1alpha1.KsqldbCluster
				curErr = r.getObject(ctx, req, &clusterInstance)
				if curErr == nil {
					config := cluster.ClusterReconcilerConfig{BaseParam: r.base, Instance: clusterInstance}
					return ctrl.Result{}, (cluster.NewClusterReconciler(config)).Start(ctx)
				}
				break
			case KSQLDB_QUERY:
				var queryInstance v1alpha1.KsqldbQuery
				curErr = r.getObject(ctx, req, &queryInstance)
				if curErr == nil {
					config := query.QueryReconcilerConfig{BaseParam: r.base, Instance: queryInstance}
					return ctrl.Result{}, (query.NewQueryReconciler(config)).Start(ctx)
				}
				break
			}
			arrErr = append(arrErr, curErr)
		}
	}

	lo.ForEach(arrErr, func(item error, index int) {
		msg := fmt.Sprintf("unable to fetch %s", r.reconOrder[index])
		r.base.Log.Error(item, msg)
	})

	return ctrl.Result{}, client.IgnoreNotFound(curErr)
}

func (r *KsqldbReconciler) getObject(ctx context.Context, req ctrl.Request, obj client.Object) error {
	err := r.base.Get(ctx, req.NamespacedName, obj)
	if !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KsqldbReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ksqldbv1alpha1.KsqldbCluster{}).
		Complete(r)
}
