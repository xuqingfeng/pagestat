package broker

import (
	"testing"

	"github.com/xuqingfeng/pagestat/vars"
)

func TestPublish(t *testing.T) {

	testConfig := Config{
		RedisUrl:      "127.0.0.1:6379",
		RedisPassword: "redis",
	}

	b := NewBroker(testConfig)
	defer b.Stop()

	_, err := b.Client.Ping().Result()
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}

	task := vars.Task{
		UUID: "76A95DFF-DB7A-446C-8C95-A041243561FD",
		Url:  "https://example.com",
		Cron: "1m",
	}
	err = b.Publish(task)
	if err != nil {
		t.Errorf("E! publish task fail %v", err)
	}
}
