package cluster

import (
	"testing"

	"github.com/softlee-io/ksqldb-operator/pkg/util/naming"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestSecret(t *testing.T) {
	t.Run("create secret", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newSecretTask()

		action, err := task.Run(ctx, conf)
		assert.Nil(t, err)
		assert.Equal(t, CREATED, action)
	})

	t.Run("update secret", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newSecretTask()
		nsName := types.NamespacedName{Namespace: conf.Instance.Namespace, Name: naming.Secret(conf.Instance.Name)}

		err := conf.BaseParam.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:        nsName.Name,
				Namespace:   nsName.Namespace,
				Labels:      conf.Instance.Labels,
				Annotations: conf.Instance.Annotations,
			},
		})
		assert.Nil(t, err)

		action, err := task.Run(ctx, conf)
		assert.Equal(t, UPDATED, action)
		assert.Nil(t, err)

		updated := &corev1.Secret{}
		err = conf.Get(ctx, nsName, updated)
		assert.Nil(t, err)
		assert.Equal(t, nsName.Name, updated.Name)
		assert.Equal(t, nsName.Namespace, updated.Namespace)
	})

	t.Run("delete foreign deployment with operator-specific labels", func(t *testing.T) {
		ctx, conf := initTestingDep(t)
		task := newSecretTask()
		nsName := types.NamespacedName{Namespace: conf.Instance.Namespace, Name: naming.Secret(conf.Instance.Name)}

		err := conf.BaseParam.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "somename",
				Namespace:   nsName.Namespace,
				Labels:      ClusterStaticLabels(),
				Annotations: conf.Instance.Annotations,
			},
		})
		assert.Nil(t, err)

		action, err := task.Run(ctx, conf)
		assert.Equal(t, DELETED, action)
		assert.Nil(t, err)
	})
}
