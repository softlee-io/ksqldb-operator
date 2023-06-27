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
	"fmt"
	"strconv"
	"strings"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/util/helper"
	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// Env names
	bootstrapServers        = "KSQL_BOOTSTRAP_SERVERS"
	listeners               = "KSQL_LISTENERS"
	serviceID               = "KSQL_KSQL_SERVICE_ID"
	sinkReplicas            = "KSQL_KSQL_SINK_REPLICAS"
	streamReplicationFactor = "KSQL_KSQL_STREAMS_REPLICATION_FACTOR"
	internalTopicReplicas   = "KSQL_KSQL_INTERNAL_TOPIC_REPLICAS"
	securityProtocol        = "KSQL_SECURITY_PROTOCOL"
	saslMechanism           = "KSQL_SASL_MECHANISM"
	saslJaasConfig          = "KSQL_SASL_JAAS_CONFIG"
	// other consts
	port            int32 = 8088
	healthcheckPath       = "/healthcheck"
)

//var

func Container(ins ksqldbv1alpha1.KsqldbCluster) corev1.Container {
	container := corev1.Container{
		Name:           "ksqldb-cluster",
		Image:          fmt.Sprintf("confluentinc/ksqldb-server:%s", ins.Spec.Version),
		Ports:          []corev1.ContainerPort{containerPort()},
		Env:            env(ins),
		ReadinessProbe: probe(10, 10, 10), //TODO: realistic value after reciliency test
		LivenessProbe:  probe(10, 60, 15), //TODO: realistic value after reciliency test
		Resources:      resource(),
	}
	//TODO: further refinement processing: validation
	return container
}

func containerPort() corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          "api",
		ContainerPort: port,
	}
}

func env(ins ksqldbv1alpha1.KsqldbCluster) []corev1.EnvVar {
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

func probe(timeout int32, period int32, terminationGrace int64) *corev1.Probe {
	return &corev1.Probe{
		TimeoutSeconds:                timeout,
		PeriodSeconds:                 period,
		TerminationGracePeriodSeconds: helper.Pointer(int64(terminationGrace)),
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: healthcheckPath,
				Port: intstr.IntOrString{IntVal: port},
			},
		},
	}
}

func resource() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{}
}
