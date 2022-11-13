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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

func newDeploymentTask() clusterTask {
	return deploymentTask{}
}

type deploymentTask struct{}

func (t deploymentTask) Name() string {
	return "deploymentTask"
}

func (t deploymentTask) Run(ctx context.Context, config ClusterReconcilerConfig) (Action, error) {
	var (
		err               error
		actionAfterApply  Action
		actionAfterDelete Action
		desired           appsv1.Deployment
	)

	desired = t.genDesired(config.Instance)
	if actionAfterApply, err = t.apply(ctx, config, desired); err != nil {
		return actionAfterApply, err
	}

	if actionAfterDelete, err = t.delete(ctx, config, desired); err != nil {
		return actionAfterDelete, err
	}

	if actionAfterDelete == NONE {
		return actionAfterApply, err
	}

	return actionAfterDelete, err
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
			Labels:      labels,
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
	}
}

func (t deploymentTask) apply(ctx context.Context, config ClusterReconcilerConfig, desired appsv1.Deployment) (Action, error) {
	existing := &appsv1.Deployment{}

	err := controllerutil.SetControllerReference(&config.Instance, &desired, config.Scheme)
	if err != nil {
		return ERROR, fmt.Errorf("failed to set controller reference: %w", err)
	}

	nsName := types.NamespacedName{Namespace: desired.Namespace, Name: desired.Name}
	err = config.Client.Get(ctx, nsName, existing)
	if err != nil && k8serrors.IsNotFound(err) {
		if err := config.BaseParam.Client.Create(ctx, &desired); err != nil {
			return ERROR, fmt.Errorf("error on deployment creation: %w", err)
		}
		config.BaseParam.Log.Info("ksqldbcluster deployment created", "name", desired.Name, "namespace", desired.Namespace)
		return CREATED, nil
	} else if err != nil {
		return ERROR, fmt.Errorf("error on getting resource: %w", err)
	}

	updated := existing.DeepCopy()
	if updated.Labels == nil {
		existing.Labels = map[string]string{}
	}
	if updated.Annotations == nil {
		existing.Annotations = map[string]string{}
	}

	updated.Spec = desired.Spec
	updated.ObjectMeta.OwnerReferences = desired.ObjectMeta.OwnerReferences

	if updated.ObjectMeta.Annotations == nil {
		updated.ObjectMeta.Annotations = map[string]string{}
	}
	if updated.ObjectMeta.Labels == nil {
		updated.ObjectMeta.Labels = map[string]string{}
	}

	for k, v := range desired.ObjectMeta.Annotations {
		updated.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desired.ObjectMeta.Labels {
		updated.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(existing)
	if err = config.Client.Patch(ctx, updated, patch); err != nil {
		return ERROR, fmt.Errorf("error on resource patch: %w", err)
	}
	config.Log.Info("applied", "deployment.name", desired.Name, "deployment.namespace", desired.Namespace)
	return UPDATED, nil
}

func (t deploymentTask) delete(ctx context.Context, config ClusterReconcilerConfig, desired appsv1.Deployment) (Action, error) {
	opts := []client.ListOption{
		client.InNamespace(config.Instance.Namespace),
		client.MatchingLabels(ClusterStaticLabels()),
	}
	list := &appsv1.DeploymentList{}
	if err := config.List(ctx, list, opts...); err != nil {
		return ERROR, fmt.Errorf("error on listing resources: %w", err)
	}
	if len(list.Items) == 0 {
		return NONE, nil
	}

	action := NONE
	for _, dep := range list.Items {
		shouldDelete := true
		if dep.Name == desired.Name && dep.Namespace == desired.Namespace {
			shouldDelete = false
		}

		if shouldDelete {
			action = DELETED
			if err := config.Delete(ctx, &dep); err != nil && !k8serrors.IsNotFound(err) {
				return ERROR, fmt.Errorf("error on deleting resource: %w", err)
			}
			config.Log.Info("deleted", "deployment.name", dep.Name, "deployment.namespace", dep.Namespace)
		}
	}

	return action, nil
}
