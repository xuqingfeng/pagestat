// Package worker is used to process tasks
package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/xuqingfeng/pagestat/trace"
	"github.com/xuqingfeng/pagestat/vars"
)

type Worker struct {
	Client *redis.Client
}

func NewWorker() *Worker {

	return &Worker{}
}

func (w *Worker) Consume(subChan chan string) error {

	pubsub := w.Client.Subscribe(vars.Channel)
	//defer pubsub.Close()

	go func() {
		for {
			msgi, err := pubsub.Receive()
			if err != nil {
				log.Println(err)
			}

			switch msg := msgi.(type) {
			case *redis.Subscription:
			// do nothing
			case *redis.Message:

				// trace and TODO: store results to redis
				var t vars.Task
				err := json.Unmarshal([]byte(msg.Payload), &t)
				if err != nil {
					log.Println(err)
				}

				ret := w.trace(t)
				retInByte, err := json.Marshal(ret)
				if err != nil {
					log.Println(err)
				}
				subChan <- string(retInByte)
			}
		}
	}()

	return nil
}

func (w *Worker) trace(task vars.Task) map[string]time.Duration {

	l := make(map[string]time.Duration)
	l["dns_lookup"] = 0
	l["tcp_connection"] = 0
	l["tls_handshake"] = 0
	l["server_processing"] = 0

	trace.Trace("GET", task.Url, []string{}, "", 0, l)

	log.Printf("I! trace result %v", l)

	return l
}

// https://www.toptal.com/go/going-real-time-with-redis-pubsub
func (w *Worker) store(uuid, result string) error {

	err := w.Client.Set(uuid, result, time.Second*120).Err()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() {

	w.Client.Close()
}
