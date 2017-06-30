package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xuqingfeng/pagestat/broker"
	"github.com/xuqingfeng/pagestat/vars"
)

type Server struct {
}

func New(c Config, b *broker.Broker) error {

	router := mux.NewRouter()
	apiRouter := router.PathPrefix(vars.APIBasePath).Subrouter()
	apiRouter.HandleFunc("/ping", PingHandler)
	// websocket connection
	apiRouter.HandleFunc("/trace", TraceHandler(b))

	if err := http.ListenAndServe(":"+strconv.Itoa(c.Port), router); err != nil {

		return err
	}

	return nil
}

func SendMessage(msg vars.Msg, status int, w http.ResponseWriter) {

	msgInByte, _ := json.Marshal(msg)
	w.WriteHeader(status)
	w.Write(msgInByte)
}
