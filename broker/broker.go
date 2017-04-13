// Package broker is used to distribute tasks
package broker

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
)

type Broker struct {
	Config   *Config
	Producer *nsq.Producer
}

func NewBroker() *Broker {

	return &Broker{}
}

func (b *Broker) Publish(task vars.Task) error {

	// ping check nsq connection
	err := b.Producer.Ping()
	if err != nil {
		return err
	}

	dataInJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return b.Producer.Publish(vars.Topic, dataInJSON)
}

func (b *Broker) Stop () {

	b.Producer.Stop()
}
