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
	"fmt"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	if err != nil && k8serrors.IsNotFound(err) {
		if err := config.BaseParam.Client.Create(ctx, &desired); err != nil {
			return fmt.Errorf("error on deployment creation: %w", err)
		}
		config.BaseParam.Log.Info("ksqldbcluster deployment created", "name", desired.Name, "namespace", desired.Namespace)
	} else if err != nil {
		return fmt.Errorf("error on getting resource: %w", err)
	}
	//TODO: validation

	//TODO:
	//- Route through ErrorHnadling (doesExist)
	//- IsNotFound error
	//	- createDesired
	//	- compareWithDesired
	//- Other Errors
	//	- returen directly
	//- IfFound
	//	- patchResource

	//TODO: HPA
	return fmt.Errorf("not implemented")
}

func (t deploymentTask) genDesired(ins ksqldbv1alpha1.KsqldbCluster) appsv1.Deployment {
	// ref: https://docs.ksqldb.io/en/latest/operate-and-deploy/installation/install-ksqldb-with-docker/
	ins.Spec.DefaultServiceID(ins.Namespace, ins.Name)
	labels := ClusterLabels(ins)
	annotations := ClusterAnnotations(ins)

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Deployment(ins.Name),
			Namespace:   ins.Namespace,
			Labels:      ins.Labels,
			Annotations: ins.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &ins.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{Container(ins)},
				},
				//TODO: NodeSelector, Affinity, Toleration
			},
		},
		//TODO:
	}
}
