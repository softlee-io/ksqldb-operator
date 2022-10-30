package config

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BaseParam struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}
