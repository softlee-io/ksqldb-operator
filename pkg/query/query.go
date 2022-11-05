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

package query

import (
	"context"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
	"github.com/softlee-io/ksqldb-operator/pkg/config"
)

type QueryReconcilerConfig struct {
	config.BaseParam
	Instance ksqldbv1alpha1.KsqldbQuery
}

func (c QueryReconcilerConfig) Validate() error {
	//TODO: Impl validation logic
	return nil
}

type queryReconciler struct {
	config QueryReconcilerConfig
}

func NewQueryReconciler(config QueryReconcilerConfig) queryReconciler {
	return queryReconciler{
		config: config,
	}
}

// No provision task included:
// - ksqldb api calls towards selected ksqldb cluster (Pod DNS)
// - desiredState = REST Get against one of cluster pods
//   - Save permanent push-based query as configmap
//   - Create: If difference is detected, send update request ()
//   - Update
//
// -
func (r queryReconciler) Start(ctx context.Context) error {
	return nil
}
