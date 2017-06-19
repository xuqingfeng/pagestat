package worker

import (
	"testing"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestConsumer(t *testing.T) {

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

	w := NewWorker()
	w.Client = testClient

	t.Log("I! comsuming")
	subChan := make(chan string)
	err = w.Consume(subChan)
	if err != nil {
		t.Errorf("E! consume task fail %v", err)
	}

	// publish message
	t.Log("I! publishing")
	pubClient := redis.NewClient(&redis.Options{
		Addr: testConfig.RedisUrl,
		Password: testConfig.RedisPassword,
	})
	task := vars.Task{
		UUID: "76A95DFF-DB7A-446C-8C95-A041243561FD",
		Url:  "https://example.com",
		Cron: "1m",
	}
	taskInByte, err := json.Marshal(task)
	if err != nil {
		t.Errorf("E! json marshal fail %v", err)
	}
	_, err = pubClient.Publish(vars.Channel, string(taskInByte)).Result()
	if err != nil {
		t.Errorf("E! redis publish message fail %v", err)
	}

	val := <-subChan
	t.Logf("I! subChan %v", val)
}
