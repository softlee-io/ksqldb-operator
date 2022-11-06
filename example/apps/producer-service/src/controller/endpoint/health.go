package endpoint

import (
	"net/http"

	"github.com/leewoobin789/test-camunda/producer-service/src/controller"
)

type healthEndpoint struct {
	info controller.HandlerInfo
}

func newhealthEndpoint() controller.Handler {
	return healthEndpoint{
		info: controller.HandlerInfo{
			Path:   "/health",
			Method: controller.GET,
		},
	}
}

func (e healthEndpoint) GetInfo() controller.HandlerInfo {
	return e.info
}

func (e healthEndpoint) Run(w http.ResponseWriter, r *http.Request) {
	controller.RespondWithJSON(w, "healthy")
}
