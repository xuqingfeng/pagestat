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
	Client    *redis.Client
	SubClient *redis.Client
}

func New(c Config) *Worker {

	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisUrl,
		Password: c.RedisPassword,
	})
	subClient := redis.NewClient(&redis.Options{
		Addr:     c.RedisUrl,
		Password: c.RedisPassword,
	})

	return &Worker{client, subClient}
}

func (w *Worker) Consume(subChan chan string) error {

	pubsub := w.SubClient.Subscribe(vars.Channel)
	//defer pubsub.Close()
	//defer close(subChan)

	type listElement struct {
		UUID string `json:"uuid"`
		Ret  string `json:"ret"`
	}

	go func() {
		for {
			msgi, err := pubsub.Receive()
			if err != nil {
				log.Println(err)
				continue
			}

			switch msg := msgi.(type) {
			case *redis.Subscription:
			// do nothing
			case *redis.Message:

				// trace
				var t vars.Task
				err := json.Unmarshal([]byte(msg.Payload), &t)
				if err != nil {
					log.Println(err)
					continue
				}

				ret := w.trace(t)
				retInByte, err := json.Marshal(ret)
				if err != nil {
					log.Println(err)
					continue
				}
				le := listElement{
					t.UUID,
					string(retInByte),
				}
				leInByte, err := json.Marshal(le)
				if err != nil {
					log.Println(err)
					continue
				}
				subChan <- string(leInByte)
			}
		}
	}()

	// store results to redis
	go func() {
		for {
			select {
			case leInString := <-subChan:

				log.Printf("I! leInString %s", leInString)
				var le listElement
				if err := json.Unmarshal([]byte(leInString), &le); err != nil {
					log.Println(err)
					continue
				}
				_, err := w.Client.LPush(le.UUID, le.Ret).Result()
				if err != nil {
					log.Println(err)
					continue
				}
			default:
				time.Sleep(time.Nanosecond * 100)
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
