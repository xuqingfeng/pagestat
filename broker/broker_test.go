package broker

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestPublish(t *testing.T) {

	testConfig := NewConfig()
	testConfig.RedisUrl = "127.0.0.1:6379"
	testConfig.RedisPassword = "redis"

	testClient := redis.NewClient(&redis.Options{
		Addr:     testConfig.RedisUrl,
		Password: testConfig.RedisPassword,
	})
	_, err := testClient.Ping().Result()
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}

	b := NewBroker()
	b.Client = testClient
	defer b.Stop()

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
