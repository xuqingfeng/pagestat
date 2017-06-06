package broker

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestPublish(t *testing.T) {

	testConfig := NewConfig()
	testConfig.RedisURL = "redis://127.0.0.1:6379"
	testConfig.RedisPassword = "redis"

	do := redis.DialPassword(testConfig.RedisPassword)
	testConn, err := redis.DialURL(testConfig.RedisURL, do)
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}

	broker := NewBroker()
	broker.Config = testConfig
	broker.Conn = testConn
	defer broker.Stop()

	task := vars.Task{
		UUID: "76A95DFF-DB7A-446C-8C95-A041243561FD",
		Url:  "https://example.com",
		Cron: "1m",
	}
	err = broker.Publish(task)
	if err != nil {
		t.Errorf("E! publish task fail %v", err)
	}
}
