package worker

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestConsumer(t *testing.T) {

	testConfig := Config{
		RedisUrl:      "127.0.0.1:6379",
		RedisPassword: "redis",
	}

	w := New(testConfig)
	_, err := w.Client.Ping().Result()
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}
	_, err = w.SubClient.Ping().Result()
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}

	t.Log("I! comsuming")
	subChan := make(chan string)
	err = w.Consume(subChan)
	if err != nil {
		t.Errorf("E! consume task fail %v", err)
	}

	// publish message
	t.Log("I! publishing")
	testPubClient := redis.NewClient(&redis.Options{
		Addr:     testConfig.RedisUrl,
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
	_, err = testPubClient.Publish(vars.Channel, string(taskInByte)).Result()
	if err != nil {
		t.Errorf("E! redis publish message fail %v", err)
	}

	time.Sleep(time.Second * 5)
}
