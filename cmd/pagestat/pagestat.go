package main

import (
	"flag"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/xuqingfeng/pagestat/broker"
	"github.com/xuqingfeng/pagestat/worker"
)

var (
	mode          string
	redisUrl      string
	redisPassword string
)

func main() {

	finish := make(chan bool)

	// broker/worker mode
	flag.StringVar(&mode, "mode", "", "mode(broker/worker)")
	flag.StringVar(&redisUrl, "redisUrl", "", "redis address")
	flag.StringVar(&redisPassword, "redisPassword", "", "redis password")
	flag.Parse()

	switch mode {
	case "broker":

		b := broker.NewBroker()
		b.Config = broker.NewConfig()
		b.Config.RedisURL = redisUrl
		b.Config.RedisPassword = redisPassword
		brokerConn, err := redis.DialURL(b.Config.RedisURL, redis.DialPassword(b.Config.RedisPassword))
		if err != nil {
			log.Fatal(err)
		}
		b.Conn = brokerConn
		defer b.Stop()

		// TODO: PUBLISH

	case "worker":

		w := worker.NewWorker()
		w.Config = worker.NewConfig()
		w.Config.RedisUrl = redisUrl
		w.Config.RedisPassword = redisPassword
		workerConn, err := redis.DialURL(w.Config.RedisUrl, redis.DialPassword(w.Config.RedisPassword))
		if err != nil {
			log.Fatal(err)
		}
		w.Conn = workerConn
		defer w.Stop()

		err = w.Consume()
		if err != nil {
			log.Printf("E! consume task fail %s", err.Error())
		}
	default:
		flag.Usage()
	}

	<-finish
}
