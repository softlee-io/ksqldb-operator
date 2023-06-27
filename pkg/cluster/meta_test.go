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

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterLabels(t *testing.T) {
	// prepare
	expected := map[string]string{
		"app.kubernetes.io/name":      "ksqldb.softlee.io/KsqldbCluster/my-cluster",
		"app.kubernetes.io/instance":  "ksqldb.softlee.io/KsqldbCluster/some-ns.my-cluster",
		"app.kubernetes.io/version":   "latest",
		"app.kubernetes.io/component": "KsqldbCluster",
		"app.kubernetes.io/part-of":   "ksqldb.softlee.io",
	}

	c := ksqldbv1alpha1.KsqldbCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "my-cluster", Namespace: "some-ns"},
		Spec:       ksqldbv1alpha1.KsqldbClusterSpec{Version: "latest"},
	}

	// test
	result := ClusterLabels(c)

	// verify
	assert.Equal(t, expected, result)
}

func TestClusterAnnotations(t *testing.T) {
	testSet := []struct {
		SaslJaasConfig string
		Expected       map[string]string
	}{
		{
			SaslJaasConfig: "",
			Expected:       map[string]string{},
		},

		{
			SaslJaasConfig: "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"<username>\" password=\"<strong-password>\";",
			Expected: map[string]string{
				"ksqldbcluster-jaas-config/sha256": "e8ae0190cff240b5592d70c0c7b47a6646cae38b1771c5886ada95bb7790ffd8",
			},
		},
	}

	for _, tc := range testSet {
		c := ksqldbv1alpha1.KsqldbCluster{
			ObjectMeta: metav1.ObjectMeta{Name: "my-cluster", Namespace: "some-ns"},
			Spec:       ksqldbv1alpha1.KsqldbClusterSpec{Version: "latest", SaslJaasConfig: tc.SaslJaasConfig},
		}
		res := ClusterAnnotations(c)
		assert.Equal(t, tc.Expected, res)
	}
}
