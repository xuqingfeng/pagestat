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

// TODO: fix hold on in test
func (w *Worker) Consume(redisUrl, redisPassword string) error {

	pubsub := w.Client.Subscribe(vars.Channel)
	//defer pubsub.Close()

	newClient := redis.NewClient(&redis.Options{
		Addr: redisUrl,
		Password: redisPassword,
	})
	//defer newClient.Close()

	for {
		msgi, err := pubsub.ReceiveTimeout(time.Second * 5)
		if err != nil {
			log.Println(err)
		}
		//
		//// trace and store results to redis
		//var t vars.Task
		//err = json.Unmarshal([]byte(msgi.Payload), &t)
		//if err != nil {
		//	return err
		//}
		//
		//ret := w.trace(t)
		//retInByte, err := json.Marshal(ret)
		//if err != nil {
		//	log.Println(err)
		//	return err
		//}

		//err = w.store(newClient, t.UUID, string(retInByte))
		//if err != nil {
		//	log.Println(err)
		//	return err
		//}

		switch msg := msgi.(type) {
		case *redis.Subscription:
		// do nothing
		case *redis.Message:

			// trace and store results to redis
			var t vars.Task
			err := json.Unmarshal([]byte(msg.Payload), &t)
			if err != nil {
				return err
			}

			ret := w.trace(t)
			retInByte, err := json.Marshal(ret)
			if err != nil {
				log.Println(err)
				return err
			}
			err = w.store(newClient, t.UUID, string(retInByte))
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
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

func (w *Worker) store(client *redis.Client, uuid, result string) error {

	client.Set("test", "test", time.Second * 120)
	err := client.Set(uuid, result, time.Second * 120).Err()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() {

	w.Client.Close()
}
