package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	"github.com/geometry-labs/go-service-template/core"
)

type KafkaTopicProducer struct {
	BrokerURL string
	TopicName string
	TopicChan chan *sarama.ProducerMessage
}

// map[Topic_Name] -> Producer
var KafkaTopicProducers = map[string]*KafkaTopicProducer{}

func StartProducers() {
	kafka_broker := core.Vars.KafkaBrokerURL
	producer_topics := core.Vars.ProducerTopics

	log.Debug("Start Producer: kafka_broker=", kafka_broker, " producer_topics=", producer_topics)

	for _, t := range producer_topics {
		KafkaTopicProducers[t] = &KafkaTopicProducer{
			kafka_broker,
			t,
			make(chan *sarama.ProducerMessage),
		}

		go KafkaTopicProducers[t].produceTopic()
	}
}

func (k *KafkaTopicProducer) produceTopic() {
	sarama_config := sarama.NewConfig()
	sarama_config.Producer.Partitioner = sarama.NewRandomPartitioner
	sarama_config.Producer.RequiredAcks = sarama.WaitForAll
	sarama_config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{k.BrokerURL}, sarama_config)
	if err != nil {
		log.Panic("KAFKA PRODUCER NEWSYNCPRODUCER PANIC: ", err.Error())
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic("KAFKA PRODUCER CLOSE PANIC: ", err.Error())
		}
	}()

	log.Debug("Producer ", k.TopicName, ": Started producing")
	for {
		topic_msg := <-k.TopicChan

		partition, offset, err := producer.SendMessage(topic_msg)
		if err != nil {
			log.Warn("Producer ", k.TopicName, ": Err sending message=", err.Error())
		}

		log.Debug("Producer ", k.TopicName, ": Producing message partition=", partition, " offset=", offset)
	}
}
