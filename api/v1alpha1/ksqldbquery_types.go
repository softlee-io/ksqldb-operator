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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type KsqldbQueryType string

const (
	ONE_TIME KsqldbQueryType = "onetime"
	STREAM   KsqldbQueryType = "stream"
)

// KsqldbQuerySpec defines the desired state of KsqldbQuery
type KsqldbQuerySpec struct {
	// Name of KsqldbCluster Resource in the same namespace
	KsqldbCluster string `json:"ksqldbCluster"`
	// KSQLDB Query
	Query string `json:"query"`
	// KSQLDB Query type
	// allowed value: "onetime", "stream" (default)
	// + optional
	QueryType KsqldbQueryType `json:"queryType,omitempty"`
}

// KsqldbQueryStatus defines the observed state of KsqldbQuery
type KsqldbQueryStatus struct {
	QueryStatus string `json:"queryStatus"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="queryStatus",type="string",JSONPath=".spec.status.queryStatus",description="Status of query"
// +operator-sdk:csv:customresourcedefinitions:displayName="KSQLDB Query"

// KsqldbQuery is the Schema for the ksqldbqueries API
type KsqldbQuery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KsqldbQuerySpec   `json:"spec,omitempty"`
	Status KsqldbQueryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KsqldbQueryList contains a list of KsqldbQuery
type KsqldbQueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KsqldbQuery `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KsqldbQuery{}, &KsqldbQueryList{})
}
