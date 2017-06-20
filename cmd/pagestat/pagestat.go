package main

import (
	"flag"
	"log"

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
	flag.StringVar(&redisUrl, "redisUrl", "", "redis url")
	flag.StringVar(&redisPassword, "redisPassword", "", "redis password")
	flag.Parse()

	switch mode {
	case "broker":

		b := broker.NewBroker(broker.Config{RedisUrl: redisUrl, RedisPassword: redisPassword})
		err := b.Client.Ping().Err()
		if err != nil {
			log.Fatalf("E! create redis connection fail %v", err)
		}

		defer b.Stop()

		// TODO: listen API request and PUBLISH message

	case "worker":

		w := worker.NewWorker(worker.Config{RedisUrl: redisUrl, RedisPassword: redisPassword})
		err := w.Client.Ping().Err()
		if err != nil {
			log.Fatalf("E! create redis connection fail %v", err)
		}
		err = w.SubClient.Ping().Err()
		if err != nil {
			log.Fatalf("E! create redis connection fail %v", err)
		}

		defer w.Stop()

		subChan := make(chan string)
		err = w.Consume(subChan)
		if err != nil {
			log.Printf("E! consume task fail %s", err.Error())
		}
	default:
		flag.Usage()
	}

	<-finish
}
