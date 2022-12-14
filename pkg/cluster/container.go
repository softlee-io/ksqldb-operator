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
	"fmt"
	"strconv"
	"strings"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/util/helper"
	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"
	corev1 "k8s.io/api/core/v1"
)

const (
	bootstrapServers        = "KSQL_BOOTSTRAP_SERVERS"
	listeners               = "KSQL_LISTENERS"
	serviceID               = "KSQL_KSQL_SERVICE_ID"
	sinkReplicas            = "KSQL_KSQL_SINK_REPLICAS"
	streamReplicationFactor = "KSQL_KSQL_STREAMS_REPLICATION_FACTOR"
	internalTopicReplicas   = "KSQL_KSQL_INTERNAL_TOPIC_REPLICAS"
	securityProtocol        = "KSQL_SECURITY_PROTOCOL"
	saslMechanism           = "KSQL_SASL_MECHANISM"
	saslJaasConfig          = "KSQL_SASL_JAAS_CONFIG"
)

func Container(ins ksqldbv1alpha1.KsqldbCluster) corev1.Container {
	container := corev1.Container{
		Name:  "ksqldb-cluster",
		Image: fmt.Sprintf("confluentinc/ksqldb-server:%s", ins.Spec.Version),
		Ports: []corev1.ContainerPort{ContainerPort()},
		Env:   Env(ins),
		//TODO: Resource, Livenessprobe, Readinessprobe
	}
	//TODO: further refinement processing: validation
	return container
}

func ContainerPort() corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          "api",
		ContainerPort: 8088,
	}
}

func Env(ins ksqldbv1alpha1.KsqldbCluster) []corev1.EnvVar {
	res := []corev1.EnvVar{
		{
			Name:  bootstrapServers,
			Value: strings.Join(ins.Spec.BootstrapServers, ","),
		},
		{
			Name:  listeners,
			Value: "http://0.0.0.0:8088/",
		},
		{
			Name:  serviceID,
			Value: ins.Spec.ServiceID,
		},
		{
			Name:  sinkReplicas,
			Value: strconv.Itoa(ins.Spec.SinkReplicas),
		},
		{
			Name:  streamReplicationFactor,
			Value: strconv.Itoa(ins.Spec.StreamReplicationFactor),
		},
		{
			Name:  internalTopicReplicas,
			Value: strconv.Itoa(ins.Spec.InternalTopicReplicas),
		},
	}

	if len(ins.Spec.SecurityProtocol) > 0 {
		res = append(res, corev1.EnvVar{
			Name:  securityProtocol,
			Value: ins.Spec.SecurityProtocol,
		})
	}

	if len(ins.Spec.SaslMechanism) > 0 {
		res = append(res, corev1.EnvVar{
			Name:  saslMechanism,
			Value: ins.Spec.SaslMechanism,
		})
	}

	if len(ins.Spec.SaslJaasConfig) > 0 {
		// addd envFrom and reference to secret
		res = append(res, corev1.EnvVar{
			Name: saslMechanism,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: naming.Secret(ins.Name),
					},
					Key:      keyJaasConfig,
					Optional: helper.Pointer(false),
				},
			},
		})
	}

	return res
}
