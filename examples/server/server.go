package main

import (
//	"fmt"
	"math/rand"
	"strconv"
	"time"
	"io/ioutil"
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	server := canopus.NewServer()

	server.Get("/hello", func(req canopus.Request) canopus.Response {
		log.Println("Hello Called")
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged with response : " + req.GetMessage().GetPayload().String())

		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Post("/hello", func(req canopus.Request) canopus.Response {
		log.Println("Hello Called via POST")

		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Get("/basic", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), canopus.NewPlainTextPayload("Acknowledged"))
		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Get("/basic/json", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), nil)

		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/watch/this", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), canopus.NewPlainTextPayload("Acknowledged"))
		res := canopus.NewResponse(msg, nil)

		return res
	})

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				//fmt.Println("[SERVER << ] Change of value -->", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()

	server.Post("/blockupload", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		// Save to file
		payload := req.GetMessage().GetPayload().GetBytes()
		log.Println("len", len(payload))
		err := ioutil.WriteFile("output.html", payload, 0644)
		if err != nil {
			log.Println(err)
		}

		return res
	})

	server.OnMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})

//	server.OnObserve(func(resource string, msg canopus.Message) {
//		fmt.Println("[SERVER << ] Observe Requested for " + resource)
//	})

//	server.OnBlockMessage(func(msg canopus.Message, inbound bool) {
//		   log.Println("Incoming Block Message:")
//		   canopus.PrintMessage(msg)
//	})

	server.ListenAndServe(":5683")
	<-make(chan struct{})
}
