package cluster

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-logr/logr/testr"
	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
	"github.com/softlee-io/ksqldb-operator/pkg/util/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func initTestingDep(t *testing.T) (context.Context, ClusterReconcilerConfig) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	log := testr.New(t)

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(ksqldbv1alpha1.AddToScheme(scheme))
	randStr := helper.RandomString(10)

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
				SecurityProtocol:        "SASL_SSL",
				SaslMechanism:           "PLAIN",
				SaslJaasConfig:          "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"admin\" password=\"admin-secret\"",
			},
		},
	}

	return ctx, conf
}
