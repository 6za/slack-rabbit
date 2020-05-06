package main

import (
	"log"
	"net/http"
	"github.com/streadway/amqp"
	"fmt"
	"encoding/json"	
	"app/slackevents"
	"os"
)
//"encoding/json"	
import (
	"io/ioutil"
 )
 import "bytes"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	//os.Exit(3)
}

type Message struct {
	Channel       string  `json:"channel"`
	Text         string   `json:"text"`
  } 


var slack_key = os.Getenv("SLACK_TOKEN")
var queue = os.Getenv("IN_QUEUE")
var queueHostname = os.Getenv("QUEUE_HOSTNAME")
var queueUser = os.Getenv("QUEUE_USER")
var queuePassword = os.Getenv("QUEUE_PASSWORD")

func main() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", queueUser,queuePassword,queueHostname))	
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var messageIn slackevents.Message
			json.Unmarshal(d.Body, &messageIn)
			//res, err := http.Get("http://www.google.com/robots.txt")
			destination := messageIn.ReplyTo.Channel
			message := messageIn.Text
			m := Message{Channel: destination, Text:message}
			b, err := json.Marshal(m)			

			client := &http.Client{}
			req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage",bytes.NewBuffer(b))
			req.Header.Add("Authorization", "Bearer " + slack_key)
			req.Header.Add("Content-Type", "application/json")		
			req.Header.Add("Charset", "utf-8")	
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			robots, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		  
			fmt.Printf("==> PING %s\n",robots)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
