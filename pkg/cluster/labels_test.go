package cluster_test

import (
	"testing"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/cluster"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSelectorLabels(t *testing.T) {
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
	result := cluster.SelectorLabels(c)

	// verify
	assert.Equal(t, expected, result)
}
