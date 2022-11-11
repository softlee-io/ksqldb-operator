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
	"math/rand"
	"strings"
	"testing"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func randomString(n int) string {
	var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}

func initDeploymentDep() (context.Context, ClusterReconcilerConfig, clusterTask) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	log := klogr.New()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(ksqldbv1alpha1.AddToScheme(scheme))
	randStr := randomString(10)

	conf := ClusterReconcilerConfig{
		BaseParam: config.BaseParam{
			Client: fake.NewFakeClient(),
			Scheme: scheme,
			Log:    log,
		},
		Instance: ksqldbv1alpha1.KsqldbCluster{
			ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("name-%s", randStr), Namespace: fmt.Sprintf("ns-%s", randStr)},
			Spec: ksqldbv1alpha1.KsqldbClusterSpec{
				Version:                 "latest",
				Replicas:                1,
				BootstrapServers:        []string{"localhost:9092"},
				SinkReplicas:            1,
				StreamReplicationFactor: 1,
				InternalTopicReplicas:   1,
			},
		},
	}
	task := newDeploymentTask()

	return ctx, conf, task
}

func TestDeployment(t *testing.T) {
	t.Run("create deployment", func(t *testing.T) {
		ctx, conf, task := initDeploymentDep()
		action, err := task.Run(ctx, conf)
		assert.Nil(t, err)
		assert.Equal(t, CREATED, action)
	})

	t.Run("update deployment", func(t *testing.T) {
		ctx, conf, task := initDeploymentDep()
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

	t.Run("update deployment", func(t *testing.T) {
		ctx, conf, task := initDeploymentDep()
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
}
