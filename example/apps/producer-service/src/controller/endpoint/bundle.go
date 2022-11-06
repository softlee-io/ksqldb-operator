package endpoint

import (
	"fmt"
	"os"

	"github.com/leewoobin789/test-camunda/producer-service/src/controller"
	"github.com/leewoobin789/test-camunda/producer-service/src/producer"
)

func ReturnBundle() []controller.Handler {
	kafkaServer := os.Getenv("KAFKA_SERVER")
	schemaRegistryUrl := os.Getenv("SCHEMA_REGISTRY_SERVER")
	fmt.Println(kafkaServer, schemaRegistryUrl)
	customProducer := producer.NewCustomKafkaProducer(kafkaServer, schemaRegistryUrl)
	return []controller.Handler{
		newSendEndpoint(customProducer),
		newhealthEndpoint(),
	}
}
