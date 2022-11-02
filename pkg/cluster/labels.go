package cluster

import (
	"fmt"

	ksqldbv1alpha1 "github.com/softlee-io/ksqldb-operator/api/v1alpha1"
)

// According to recommended labels: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func SelectorLabels(ins ksqldbv1alpha1.KsqldbCluster) map[string]string {
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

	return map[string]string{
		"app.kubernetes.io/name":      OperatorPrefix(ins.Name),
		"app.kubernetes.io/instance":  OperatorPrefix(ins.Namespace, ins.Name),
		"app.kubernetes.io/version":   ins.Spec.Version,
		"app.kubernetes.io/component": "KsqldbCluster",
		"app.kubernetes.io/part-of":   "ksqldb.softlee.io",
	}
}
