// Copyright © 2018 Andrea Funtò - released under the MIT License

package rmq

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dihedron/go-log"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var (
	//ServerURL is the URL of the AMQP server
	ServerURL string
	// QueueName is the name of the RabbitmMQ queue.
	QueueName string
)

// Publish sends messages to the given AMQP server and queue.
func Publish(cmd *cobra.Command, args []string) {
	log.Infof("publishing to %q, queue %q\n", ServerURL, QueueName)

	conn, err := amqp.Dial(ServerURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Debugln("connection established")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	log.Debugln("channel opened")

	q, err := ch.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	log.Debugln("queue declared")

	forever := make(chan bool)

	go func() {
		log.Debugln("starting to publish messages...")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			var buffer strings.Builder
			var count int
			var body string
			fmt.Println("Please enter message; to exit press CTRL+C:")
			for scanner.Scan() {
				//log.Debugf("count is %d, text is %q", count, scanner.Text())
				if scanner.Text() == "" {
					if count == 0 {
						count++
						continue
					} else {
						body = buffer.String()
						break
					}
				}
				buffer.WriteString(scanner.Text())
				buffer.WriteRune('\n')
				if err := scanner.Err(); err != nil {
					log.Warnf("error reading text: %v", err)
				}
			}
			//log.Debugf("publishing message:\n%s\n", body)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			if err != nil {
				log.Errorf("failed to publish a message: %v", err)
			}
			//log.Debugf("message sent:\n%s\n", body)
		}
	}()
	<-forever
}

// Subscribe receives messages from the given AMQP server and queue.
func Subscribe(cmd *cobra.Command, args []string) {
	log.Infof("subscribed to queue %q on server %q\n", QueueName, ServerURL)

	conn, err := amqp.Dial(ServerURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Debugln("connection established")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	log.Debugln("channel opened")

	q, err := ch.QueueDeclare(
		QueueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	log.Debugln("queue declared")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		fmt.Println("Waiting for messages; to exit press CTRL+C.")
		fmt.Println("--------------------------------------------------------------------------------")
		for d := range msgs {
			fmt.Printf("%s--------------------------------------------------------------------------------\n", d.Body)
		}
	}()

	<-forever

}
