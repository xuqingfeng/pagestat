// Package broker is used to distribute tasks
package broker

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/xuqingfeng/pagestat/vars"
)

type Broker struct {
	Config *Config
	Conn   redis.Conn
}

func NewBroker() *Broker {

	return &Broker{}
}

func (b *Broker) Publish(task vars.Task) error {

	taskInBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = b.Conn.Do("PUBLISH", vars.Channel, taskInBytes)
	return err
}

func (b *Broker) Stop() {

	b.Conn.Close()
}
