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
	"testing"

	"github.com/go-logr/logr"
	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDeployment(t *testing.T) {
	ct := newDeploymentTask()
	ctx := context.Background()
	scheme := runtime.NewScheme()
	log, err := logr.FromContext(ctx)

	assert.NotNil(t, err)

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(ksqldbv1alpha1.AddToScheme(scheme))

	conf := ClusterReconcilerConfig{
		BaseParam: config.BaseParam{
			Client: fake.NewFakeClient(),
			Scheme: scheme,
			Log:    log,
		},
		//TODO:
		Instance: ksqldbv1alpha1.KsqldbCluster{
			ObjectMeta: metav1.ObjectMeta{Name: "my-cluster", Namespace: "some-ns"},
			Spec: ksqldbv1alpha1.KsqldbClusterSpec{
				Version: "latest",
			},
		},
	}

	ct.Run(ctx, conf)
}
