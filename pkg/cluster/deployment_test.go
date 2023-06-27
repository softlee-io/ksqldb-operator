// Copyright 2023 Softlee.io Authors
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
	"testing"

	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestDeployment(t *testing.T) {

	t.Run("create deployment", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newDeploymentTask()

		action, err := task.Run(ctx, conf)
		assert.Nil(t, err)
		assert.Equal(t, CREATED, action)
	})

	t.Run("update deployment", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newDeploymentTask()
		nsName := types.NamespacedName{Namespace: conf.Instance.Namespace, Name: naming.Deployment(conf.Instance.Name)}

		err := conf.BaseParam.Create(ctx, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        nsName.Name,
				Namespace:   nsName.Namespace,
				Labels:      conf.Instance.Labels,
				Annotations: conf.Instance.Annotations,
			},
		})
		assert.Nil(t, err)

		action, err := task.Run(ctx, conf)
		assert.Equal(t, UPDATED, action)
		assert.Nil(t, err)

		updated := &appsv1.Deployment{}
		err = conf.Get(ctx, nsName, updated)
		assert.Nil(t, err)
		assert.Equal(t, nsName.Name, updated.Name)
		assert.Equal(t, nsName.Namespace, updated.Namespace)
	})

	t.Run("delete foreign deployment with operator-specific labels", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newDeploymentTask()
		nsName := types.NamespacedName{Namespace: conf.Instance.Namespace, Name: naming.Deployment(conf.Instance.Name)}

		err := conf.BaseParam.Create(ctx, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "somename",
				Namespace:   nsName.Namespace,
				Labels:      ClusterStaticLabels(),
				Annotations: conf.Instance.Annotations,
			},
		})
		assert.Nil(t, err)

		action, err := task.Run(ctx, conf)
		assert.Equal(t, DELETED, action)
		assert.Nil(t, err)
	})
}
