// Package worker is used to process tasks
package worker

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
)

type Worker struct {
	Config   *Config
	Consumer *nsq.Consumer
	sync.RWMutex
}

func NewWorker() *Worker {

	return &Worker{}
}

func (w *Worker) Consume() error {

	w.Consumer.AddHandler(w)
	err := w.Consumer.ConnectToNSQLookupd(w.Config.NsqLookupdAddr)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) HandleMessage(message *nsq.Message) error {

	w.Lock()
	var task vars.Task
	err := json.Unmarshal(message.Body, &task)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile("test_data.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Printf("I! got message %+v\n", task)
	w.Unlock()
	return nil
}

func (w *Worker) Stop() {

	w.Consumer.Stop()
}
