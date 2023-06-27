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
	"crypto/sha256"
	"fmt"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
)

// According to recommended labels: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func ClusterLabels(ins ksqldbv1alpha1.KsqldbCluster) map[string]string {
	OperatorPrefix := func(ss ...string) string {
		var res string = ss[0]

		if len(ss) > 1 {
			for i, s := range ss {
				if i != 0 {
					res = fmt.Sprintf("%s.%s", res, s)
				}
			}
		}
		return fmt.Sprintf("ksqldb.softlee.io/KsqldbCluster/%s", res)
	}

	labels := ClusterStaticLabels()
	labels["app.kubernetes.io/name"] = OperatorPrefix(ins.Name)
	labels["app.kubernetes.io/instance"] = OperatorPrefix(ins.Namespace, ins.Name)
	labels["app.kubernetes.io/version"] = ins.Spec.Version

	if ins.Labels != nil {
		for lk, lv := range ins.Labels {
			labels[lk] = lv
		}
	}

	return labels
}

func ClusterStaticLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/component": "KsqldbCluster",
		"app.kubernetes.io/part-of":   "ksqldb.softlee.io",
	}
}

func ClusterAnnotations(ins ksqldbv1alpha1.KsqldbCluster) map[string]string {
	if len(ins.Spec.SaslJaasConfig) == 0 {
		return map[string]string{}
	}

	jaasConf := getSHA(ins.Spec.SaslJaasConfig)
	return map[string]string{
		//For purpose of integrity check for jaas auth config
		"ksqldbcluster-jaas-config/sha256": jaasConf,
	}
}

func getSHA(val string) string {
	h := sha256.Sum256([]byte(val))
	return fmt.Sprintf("%x", h)
}
