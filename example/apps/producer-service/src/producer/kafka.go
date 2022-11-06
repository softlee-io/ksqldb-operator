package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
)

type CustomProducer interface {
	Send(topic string, key string, value interface{}) error
}

type CustomKafkaProducer struct {
	producer   *kafka.Producer
	serializer *avro.SpecificSerializer
}

func NewCustomKafkaProducer(kafkaServer string, schemaRegistryUrl string) CustomProducer {
	config := schemaregistry.NewConfig(schemaRegistryUrl)
	client, err := schemaregistry.NewClient(config)

	if err != nil {
		panic(err)
	}

	avroSerializer, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())

	if err != nil {
		panic(err)
	}

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		panic(err)
	}
	return CustomKafkaProducer{
		producer:   kafkaProducer,
		serializer: avroSerializer,
	}
}

func (k CustomKafkaProducer) Send(topic string, key string, value interface{}) error {
	payload, err := k.serializer.Serialize(topic, value)
	if err != nil {
		return err
	}

	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          payload,
	}, nil)
}
