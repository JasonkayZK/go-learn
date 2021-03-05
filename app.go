package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/events", es)
	go func() {
		for {
			// 每2秒发送一条当前时间消息，并打印对应客户端数量
			es.SendEventMessage(fmt.Sprintf("hello, now is: %s", time.Now()), "", "")
			log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(2 * time.Second)
		}
	}()

	log.Println("Open URL http://localhost:8080/ in your browser.")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
