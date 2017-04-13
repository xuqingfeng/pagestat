package main

import (
	"flag"
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/nsqio/go-nsq"
	"github.com/xuqingfeng/pagestat/vars"
	"github.com/xuqingfeng/pagestat/worker"
)

var (
	mode       string
	nsqlookupd string
)

func main() {

	finish := make(chan bool)

	// broker/worker mode
	flag.StringVar(&mode, "mode", "", "mode(broker/workder)")
	flag.StringVar(&nsqlookupd, "nsqlookupd", "", "nsqlookupd http address")
	flag.Parse()

	channelName, err := os.Hostname()
	if err != nil {
		channelName = "undefined"
	}

	switch mode {
	case "broker":

	case "worker":
		consumer, err := nsq.NewConsumer(vars.Topic, channelName, nsq.NewConfig())
		if err != nil {
			log.Fatalf("E! create nsq consumer fail %s", err.Error())
		}
		w := worker.NewWorker()
		w.Config = worker.NewConfig()
		w.Config.NsqLookupdAddr = nsqlookupd
		w.Consumer = consumer
		defer w.Stop()

		err = w.Consume()
		if err != nil {
			log.Errorf("E! consume task fail %s", err.Error())
		}
	default:
		flag.Usage()
	}

	<-finish
}
