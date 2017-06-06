package worker

import (
	"log"
	"os"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
)

func TestConsumer(t *testing.T) {

	testConfig := NewConfig()
	testConfig.RedisUrl = "redis://127.0.0.1:6379"
	testConfig.RedisPassword = "redis"

	do := redis.DialPassword(testConfig.RedisPassword)
	testConn, err := redis.DialURL(testConfig.RedisUrl, do)
	if err != nil {
		t.Fatalf("E! create redis connection fail %v", err)
	}

	worker := NewWorker()
	worker.Config = testConfig
	worker.Conn = testConn
	defer worker.Stop()

	err = worker.Consume()
	if err != nil {
		t.Errorf("E! consume task fail %v", err)
	}
}
