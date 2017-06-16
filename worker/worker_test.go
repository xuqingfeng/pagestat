package worker

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/xuqingfeng/pagestat/broker"
	"github.com/xuqingfeng/pagestat/vars"
	"time"
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

	b := broker.NewBroker()
	b.Client = testClient

	//defer testClient.Close()

	//err = testClient.Set("test", "test", time.Second * 200).Err()
	//if err != nil {
	//	t.Errorf("E! set fail %v", err)
	//}

	go func() {
		t.Log("I! comsuming")

		err = w.Consume(testConfig.RedisUrl, testConfig.RedisPassword)
		if err != nil {
			t.Errorf("E! consume task fail %v", err)
		}
	}()

	task := vars.Task{
		UUID: "76A95DFF-DB7A-446C-8C95-A041243561FD",
		Url:  "https://example.com",
		Cron: "1m",
	}
	t.Log("I! publishing")
	err = b.Publish(task)
	if err != nil {
		t.Errorf("E! publish task fail %v", err)
	}

	time.Sleep(10 * time.Second)

	//select {
	//case ret := <-subChan:
	//	t.Logf("I! ret %s", ret)
	//	//var latency trace.Latency
	//	//err = json.Unmarshal(ret, &latency)
	//	//if err != nil {
	//	//	t.Errorf("E! json unmarshal fail %v", err)
	//	//}
	//	//if latency["dns_lookup"] == 0 {
	//	//	t.Error("E! dns lookup takes 0 ms")
	//	//}
	//default:
	//	t.Error("E! subChan is empty")
	//}
}
