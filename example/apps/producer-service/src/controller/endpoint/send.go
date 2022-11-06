package endpoint

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/leewoobin789/test-camunda/producer-service/src/avroschema"
	"github.com/leewoobin789/test-camunda/producer-service/src/controller"
	"github.com/leewoobin789/test-camunda/producer-service/src/producer"
)

var TOPIC = os.Getenv("KAFKA_TOPIC")

type sendEndpoint struct {
	info     controller.HandlerInfo
	producer producer.CustomProducer
}

func newSendEndpoint(producer producer.CustomProducer) controller.Handler {
	return sendEndpoint{
		info: controller.HandlerInfo{
			Path:   "/send",
			Method: controller.GET,
		},
		producer: producer,
	}
}

func (e sendEndpoint) GetInfo() controller.HandlerInfo {
	return e.info
}

func (e sendEndpoint) Run(w http.ResponseWriter, r *http.Request) {
	Response := struct {
		Message string `json:"message"`
	}{
		Message: "successful",
	}

	strNum := r.URL.Query().Get("number")
	if len(strNum) == 0 {
		strNum = "1"
	}
	num, err := strconv.Atoi(strNum)
	if err != nil {
		Response.Message = err.Error()
		controller.RespondwithJSON(w, http.StatusNotAcceptable, Response)
	}

	keyProductID := uuid.NewString()
	value := &avroschema.Order{
		Name:       "Woobin",
		FamilyName: "Lee",
		Birth:      0, // TODO:
		CustomId:   "mycustomid-1234",
		UnitPrice:  12.95,
		Amount:     5,
		Credit:     10,
		Distance:   10,
	}

	for i := 0; i <= num; i++ {
		time.Sleep(time.Millisecond * 100)
		if err := e.producer.Send(TOPIC, keyProductID, value); err != nil {
			Response.Message = err.Error()
			break
		}
	}

	controller.RespondWithJSON(w, Response)
}
