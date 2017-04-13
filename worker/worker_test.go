package worker

import (
	"log"
	"os"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestConsumer(t *testing.T) {

	testConfig := NewConfig()
	testConfig.NsqLookupdAddr = "127.0.0.1:4161"
	channelName, err := os.Hostname()
	if err != nil {
		channelName = "undefined"
	}
	consumer, err := nsq.NewConsumer(vars.Topic, channelName, nsq.NewConfig())

	logger := log.New(os.Stdout, "[pagestat] worker ", 1)
	consumer.SetLogger(logger, 1)

	if err != nil {
		t.Fatalf("E! create nsq consumer fail %v", err)
	}
	worker := NewWorker()
	worker.Config = testConfig
	worker.Consumer = consumer
	defer worker.Stop()

	err = worker.Consume()
	if err != nil {
		t.Errorf("E! consume task fail %v", err)
	}
}
