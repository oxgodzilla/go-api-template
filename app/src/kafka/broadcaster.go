package kafka

import (
	confluent "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// TODO use uuid for larger ID range
type BroadcasterID int

var LAST_BROADCASTER_ID BroadcasterID = 0

type TopicBroadcaster struct {

	// Input
	InputChan chan *confluent.Message

	// Output
	OutputChans map[BroadcasterID]chan *confluent.Message
}

var Broadcasters = map[string]*TopicBroadcaster{}

func newBroadcaster(topic_name string, input_chan chan *confluent.Message) {
	Broadcasters[topic_name] = &TopicBroadcaster{
		input_chan,
		make(map[BroadcasterID]chan *confluent.Message),
	}
}

func (tb *TopicBroadcaster) AddOutputChannel(topic_chan chan *confluent.Message) BroadcasterID {
	id := LAST_BROADCASTER_ID
	LAST_BROADCASTER_ID++

	tb.OutputChans[id] = topic_chan

	return id
}

func (tb *TopicBroadcaster) RemoveOutputChannel(id BroadcasterID) {
	_, ok := tb.OutputChans[id]
	if ok {
		delete(tb.OutputChans, id)
	}
}

func (tb *TopicBroadcaster) Broadcast() {
	for {
		msg := <-tb.InputChan

		for _, channel := range tb.OutputChans {
			select {
			case channel <- msg:
				continue
			default:
				continue
			}
		}
	}
}