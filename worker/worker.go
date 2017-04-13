// Package worker is used to process tasks
package worker

import (
	"encoding/json"
	"sync"
	"log"
	"os"

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

	w.Consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {

		var task vars.Task
		err := json.Unmarshal(m.Body, &task)
		if err != nil {
			return err
		}
		//logFile, err := os.OpenFile("test_data/worker.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		//if err != nil {
		//	return err
		//}
		//defer logFile.Close()

		//log.SetOutput(logFile)
		//log.Printf("I! got message %+v\n", task)

		os.Stdout.WriteString(task.Url)

		return nil
	}))
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
	logFile, err := os.OpenFile("test_data/worker.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Printf("I! got message %+v\n", task)
	w.Unlock()
	return nil
}

func (w *Worker) Stop(){

	w.Consumer.Stop()
}
