// Copyright 2022 Softlee.io Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"errors"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
)

type clusterTask interface {
	Run(ctx context.Context, config ClusterReconcilerConfig) error
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
	tasks  []clusterTask
}

func NewClusterReconciler(config ClusterReconcilerConfig) clusterReconciler {
	return clusterReconciler{
		config: config,
		tasks:  nil,
	}
}

func (r *clusterReconciler) initTasks() error {
	var (
		err error
	)

	err = r.config.Validate()
	if err != nil {
		return errors.New("") // TODO: precise error consolidation
	}
	// Sequential Reconciliation
	tasks := []clusterTask{
		newSecretTask(),
		newDeploymentTask(),
	}
	if !r.config.Instance.Spec.ServiceResourceDisabled {
		tasks = append(tasks, newServiceTask())
	}

	r.tasks = tasks

	return err
}

func (r clusterReconciler) Start(ctx context.Context) error {
	/*
		For each kubernetes resource
		1. Init & interate task execution
		2. desiredResource
		- Must be sequentially executed: Secret -> Deployment -> (opt. justification: ksql api can be called over pod DNS ) Service
		- for instance: constructing deployment resource configuration
		3. Error Handling
	*/
	var (
		err error
	)

	r.initTasks()
	for _, task := range r.tasks {
		err = task.Run(ctx, r.config)
		if err != nil {
			return err
		}
	}

	return err
}
