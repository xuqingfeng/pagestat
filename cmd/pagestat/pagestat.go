package main

import (
	"flag"
	"log"
	"os"

	"github.com/xuqingfeng/pagestat/broker"
	"github.com/xuqingfeng/pagestat/server"
	"github.com/xuqingfeng/pagestat/worker"
)

var (
	mode          string
	redisUrl      string
	redisPassword string
	serverPort    int
)

func main() {

	finish := make(chan bool)

	// broker/worker mode
	flag.StringVar(&mode, "mode", "", "mode (broker/worker)")
	flag.StringVar(&redisUrl, "redis-url", "127.0.0.1:6379", "redis url")
	flag.StringVar(&redisPassword, "redis-password", "redis", "redis password")
	flag.IntVar(&serverPort, "server-port", 2017, "server port")
	flag.Parse()

	switch mode {
	case "broker":

		b := broker.New(broker.Config{RedisUrl: redisUrl, RedisPassword: redisPassword})
		err := b.Client.Ping().Err()
		if err != nil {
			log.Fatalf("E! create redis connection fail %v", err)
		}

		defer b.Stop()

		// TODO: listen API request and PUBLISH message
		err = server.New(server.Config{Port: serverPort}, b)
		if err != nil {
			log.Fatalf("E! create server fail %v", err)
		}

	case "worker":

		w := worker.New(worker.Config{RedisUrl: redisUrl, RedisPassword: redisPassword})
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
		os.Exit(0)
	}

	<-finish
}
