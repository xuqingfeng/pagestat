package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/xuqingfeng/pagestat/broker"
	"github.com/xuqingfeng/pagestat/vars"
)

type data struct {
	Url string `json:"url"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func PingHandler(w http.ResponseWriter, r *http.Request) {

	msg := vars.Msg{
		Success: true,
		Message: "pong",
	}

	SendMessage(msg, http.StatusOK, w)
}

func TraceHandler(b *broker.Broker) http.HandlerFunc {

	id := uuid.NewV4().String()

	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			msg := vars.Msg{
				Success: false,
				Message: err.Error(),
			}
			log.Printf("E! upgrade fail %v", err)
			msgInByte, _ := json.Marshal(msg)
			w.Write(msgInByte)
			return
		}

		send := make(chan vars.Msg)
		ticker := time.NewTicker(time.Second * 5)

		// read message
		go func() {
			defer conn.Close()
			m := vars.Msg{
				Success: false,
			}
			for {
				msgType, msg, err := conn.ReadMessage()
				log.Println("I! read message")
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
						ticker.Stop()
						log.Printf("E! read message fail %v", err)
					}
					m.Message = err.Error()
					send <- m
					break
				}
				if msgType == websocket.TextMessage {
					log.Printf("I! get text message")
					var d data
					err = json.Unmarshal(msg, &d)
					if err != nil {
						m.Message = err.Error()
						send <- m
						break
					}
					t := vars.Task{
						UUID: id,
						Url:  d.Url,
						Cron: "10m",
					}
					err = b.Publish(t)
					if err != nil {
						m.Message = err.Error()
						send <- m
						break
					}
				}

			}
		}()

		// write message
		go func() {
			defer conn.Close()
			for {
				select {
				case m := <-send:
					mInByte, _ := json.Marshal(m)
					conn.WriteMessage(websocket.TextMessage, mInByte)
				case <-ticker.C:

					rets := b.Client.LRange(id, 0, -1).Val()
					m := vars.Msg{
						Success: true,
						Data:    rets,
						Message: "get results",
					}
					log.Printf("I! get results %v", rets)
					mInByte, _ := json.Marshal(m)
					// return when connection close
					if err := conn.WriteMessage(websocket.TextMessage, mInByte); err != nil {
						return
					}
				}
			}
		}()
	}
}
