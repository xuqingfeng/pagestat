// Package worker is used to process tasks
package worker

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/xuqingfeng/pagestat/trace"
	"github.com/xuqingfeng/pagestat/vars"
)

type Worker struct {
	Config *Config
	Conn   redis.Conn
}

func NewWorker() *Worker {

	return &Worker{}
}

func (w *Worker) Consume() error {

	psc := redis.PubSubConn{Conn: w.Conn}
	psc.Subscribe(vars.Channel)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			go func() {
				// trace and store results to redis
				taskInBytes := byte(v.Data)
				var t vars.Task
				err := json.Unmarshal(taskInBytes, &t)
				if err != nil {
					log.Println(err)
				}
				ret := w.trace(t)
				retInBytes, err := json.Marshal(ret)
				if err != nil {
					log.Println(err)
				}
				err = w.store(t.UUID, string(retInBytes))
				if err != nil {
					log.Println(err)
				}
			}()
		case error:
			log.Println(v.Error())
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

func (w *Worker) store(uuid, result string) error {

	val, err := w.Conn.Do("SET", uuid, result, "NX", "EX", "120")
	if err != nil {
		return err
	}
	if val == nil {
		return errors.New("uuid conflict")
	}

	return nil
}

func (w *Worker) Stop() {

	w.Conn.Close()
}
