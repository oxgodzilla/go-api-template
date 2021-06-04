package kafka

import (
	"github.com/geometry-labs/go-service-template/core"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
	"time"
)

func StartApiConsumers() {
	kafka_broker := core.Vars.KafkaBrokerURL
	consumer_topics := core.Vars.ConsumerTopics

	log.Debug("Start Consumer: kafka_broker=", kafka_broker, " consumer_topics=", consumer_topics)

	BroadcastFunc := func(channel chan *sarama.ConsumerMessage, message *sarama.ConsumerMessage) {
		select {
		case channel <- message:
			return
		default:
			return
		}
	}

	for _, t := range consumer_topics {
		// Broadcaster indexed in Broadcasters map
		// Starts go routine
		newBroadcaster(t, BroadcastFunc)

		topic_consumer := &KafkaTopicConsumer{
			kafka_broker,
			t,
			Broadcasters[t],
		}

		// One routine per topic
		log.Debug("Start Consumers: Starting ", t, " consumer...")
		go topic_consumer.consumeTopic()
	}
}

func StartWorkerConsumers() {
	kafka_broker := core.Vars.KafkaBrokerURL
	consumer_topics := core.Vars.ConsumerTopics

	log.Debug("Start Consumer: kafka_broker=", kafka_broker, " consumer_topics=", consumer_topics)

	for _, t := range consumer_topics {
		// Broadcaster indexed in Broadcasters map
		// Starts go routine
		newBroadcaster(t, nil)

		topic_consumer := &KafkaTopicConsumer{
			kafka_broker,
			t,
			Broadcasters[t],
		}

		// One routine per topic
		log.Debug("Start Consumers: Starting ", t, " consumer...")
		go topic_consumer.consumeTopic()
	}
}

type KafkaTopicConsumer struct {
	BrokerURL   string
	TopicName   string
	Broadcaster *TopicBroadcaster
}

func (k *KafkaTopicConsumer) consumeTopic() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{k.BrokerURL}, config)
	if err != nil {
		log.Panic("KAFKA CONSUMER NEWCONSUMER PANIC: ", err.Error())
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Panic("KAFKA CONSUMER CLOSE PANIC: ", err.Error())
		}
	}()

	offset := sarama.OffsetOldest
	partitions, err := consumer.Partitions(k.TopicName)
	if err != nil {
		log.Panic("KAFKA CONSUMER PARTITIONS PANIC: ", err.Error(), " Topic: ", k.TopicName)
	}

	log.Debug("Consumer ", k.TopicName, ": Started consuming")
	for _, p := range partitions {
		pc, err := consumer.ConsumePartition(k.TopicName, p, offset)

		if err != nil {
			log.Panic("KAFKA CONSUMER PARTITIONS PANIC: ", err.Error())
		}
		if pc == nil {
			log.Panic("KAFKA CONSUMER PARTITIONS PANIC: Failed to create PartitionConsumer")
		}

		// Watch errors
		// go func() {
		// 		for err := range pc.Errors() {
		// log.Warn("KAFKA CONSUMER WARN: ", err.Error())
		// 	}
		// }()

		// One routine per partition
		go func(pc sarama.PartitionConsumer) {
			for {
				var topic_msg *sarama.ConsumerMessage
				select {
				case msg := <-pc.Messages():
					topic_msg = msg
				case consumerErr := <-pc.Errors():
					log.Warn("KAFKA PARTITION CONSUMER ERROR:", consumerErr.Err)
					continue
				case <-time.After(5 * time.Second):
					log.Debug("Consumer ", k.TopicName, ": No new kafka messages, waited 5 secs")
					continue
				}
				//topic_msg := <-pc.Messages()
				log.Debug("Consumer ", k.TopicName, ": Consumed message key=", string(topic_msg.Key))

				// Broadcast
				k.Broadcaster.ConsumerChan <- topic_msg

				log.Debug("Consumer ", k.TopicName, ": Broadcasted message key=", string(topic_msg.Key))
			}
		}(pc)
	}
	// Waiting, so that client remains alive
	ch := make(chan int, 1)
	<-ch
}
