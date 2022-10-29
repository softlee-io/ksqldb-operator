package cluster

import (
	"context"
	"errors"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/internal/internalreconciler/config"
)

type clusterTask interface {
	run(clusterRecon ClusterReconcilerConfig) error
}

type ClusterReconcilerConfig struct {
	config.BaseParam
	Instance ksqldbv1alpha1.KsqldbCluster
}

func (c ClusterReconcilerConfig) Validate() error {
	//TODO: Impl validation logic
	return nil
}

type clusterReconciler struct {
	config ClusterReconcilerConfig
}

func NewClusterReconciler(config ClusterReconcilerConfig) clusterReconciler {
	return clusterReconciler{
		config: config,
	}
}

func (r clusterReconciler) initTasks() error {
	err := r.config.Validate()
	if err != nil {
		return errors.New("") // TODO: precise error consolidation
	}

	//Task definition

	return err
}

func (r clusterReconciler) Start(ctx context.Context) error {
	/*
		For each kubernetes resource
		1. desiredResource
		- Must be sequentially executed: Configmap -> Secret -> Deployment -> (opt. justification: ksql api can be called over pod DNS ) Service
		- for instance: constructing deployment resource configuration
		2.
		- controllerutil.SetControllerReference(
		3.
		-
	*/
	// desiredState
	return nil
}

func (r clusterReconciler) deployment(ctx context.Context) error {
	//desired :=

	return nil
}
