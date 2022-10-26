/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// serviceID ("KSQL_KSQL_SERVICE_ID") will be created by f
type KsqldbClusterSpec struct {
	// Version/image tag of Ksqldb
	// Default: "latest"
	// +optional
	Version string `json:"version,omitempty"`
	// replicas of Ksqldb Cluster deployment
	// Default: "1"
	// +optional
	Replicas int `json:"replicas,omitempty"`
	// Bootstrap servers
	// +kubebuilder:validation:MinItems=1
	BootstrapServers []string `json:"bootstrapServers"`
	// e.g. "something_" will be used as a prefix for internal topics
	// Default: "{namespace}_{ksqldb resource name}_"
	// +optional
	ServiceID string `json:"serviceID,omitempty"`
	// number of replicas for KSQL topics (default: 1)
	// +optional
	SinkReplicas int `json:"sinkReplicas,omitempty"`
	// replication factor of KSQL internal, command and output topics (default: 3)
	// +optional
	StreamReplicationFactor int `json:"streamReplicationFactor,omitempty"`
	// replication factor of KSQL internal topics (default: 3)
	// +optional
	InternalTopicReplicas int `json:"internalTopicReplicas,omitempty"`
	// Security Protocol of KSQL Cluster (e.g. "SASL_SSL")
	// +optional
	SecurityProtocol string `json:"securityProtocol,omitempty"`
	// Sasl Mechanism used for Authentication (e.g. "PLAIN")
	// +optional
	SaslMechanism string `json:"saslMechanism,omitempty"`
	// String format of JAAS Config (e.g. "org.apache.kafka.common.security.plain.PlainLoginModule required....")
	// +optional
	SaslJaasConfig string `json:"saslJaasConfig,omitempty"`
}

// KsqldbClusterStatus defines the observed state of KsqldbCluster
type KsqldbClusterStatus struct {
	// +optional
	Version string `json:"version"`
	// +optional
	Replicas int `json:"replica"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:name="Replicas",type="int",JSONPath=".spec.replicas",description="Replicas"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version",description="KSQLDB Version"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +operator-sdk:csv:customresourcedefinitions:displayName="KSQLDB Cluster"

// KsqldbCluster is the Schema for the ksqldbclusters API
type KsqldbCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KsqldbClusterSpec   `json:"spec,omitempty"`
	Status KsqldbClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KsqldbClusterList contains a list of KsqldbCluster
type KsqldbClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KsqldbCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KsqldbCluster{}, &KsqldbClusterList{})
}
