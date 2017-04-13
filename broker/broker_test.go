package broker

import (
	"log"
	"os"
	"testing"

	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestPublish(t *testing.T) {

	testConfig := NewConfig()
	testConfig.NsqdAddr = "127.0.0.1:4150"
	producer, err := nsq.NewProducer(testConfig.NsqdAddr, nsq.NewConfig())

	logger := log.New(os.Stdout, "[pagestat] broker ", 1)
	producer.SetLogger(logger, 1)

	if err != nil {
		t.Fatalf("E! create nsq producer fail %v", err)
	}
	broker := NewBroker()
	broker.Config = testConfig
	broker.Producer = producer
	defer broker.Stop()

	task := vars.Task{
		Url:  "https://example.com",
		Cron: "1m",
	}
	err = broker.Publish(task)
	if err != nil {
		t.Errorf("E! publish task fail %v", err)
	}
}
