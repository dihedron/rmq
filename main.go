// Copyright © 2018 Andrea Funtò - released under the MIT License

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dihedron/go-log"
	"github.com/streadway/amqp"
)

func init() {
	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout, true)
	log.SetTimeFormat("15:04:05.000")
	log.SetPrintCallerInfo(true)
	log.SetPrintSourceInfo(log.SourceInfoShort)
}

func main() {
	server := flag.String("server", "", "the address of the server, including the vhost (e.g. amqp://user:passwd@server:5672/vhost)")
	queue := flag.String("queue", "", "the name of the queue to connect to")
	flag.Parse()

	log.Infof("subscribed to queue %q on server %q\n", *queue, *server)

	conn, err := amqp.Dial(*server)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Debugf("connection open:\n%s\n", GetInfoFromConnection(conn))

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	log.Debugln("channel open")

	q, err := ch.QueueDeclare(
		*queue, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	log.Debugf("queue declared:\n%s\n", GetInfoFromQueue(q))

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

// ConnectionInfo containes all metadata about a connection.
type ConnectionInfo struct {
	Vhost      string                 `json:"vhost,omitempty" yaml:"vhost,omitempty"`
	Major      int                    `json:"major" yaml:"major"`
	Minor      int                    `json:"minor" yaml:"minor"`
	Locales    []string               `json:"locales,omitempty" yaml:"locales,omitempty"`
	Server     map[string]interface{} `json:"server,omitempty" yaml:"server,omitempty"`
	ChannelMax int                    `json:"channelmax,omitempty" yaml:"channelmax,omitempty"`
	FrameSize  int                    `json:"framesize,omitempty" yaml:"framesize,omitempty"`
	Heartbeat  time.Duration          `json:"heartbeat,omitempty" yaml:"heartbeat,omitempty"`
	Locale     string                 `json:"locale,omitempty" yaml:"locale,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty"`
	TLS        *tls.Config            `json:"tls,omitempty" yaml:"tls,omitempty"`
}

// GetInfoFromConnection returns a ConnectionInfo struct for a given connection.
func GetInfoFromConnection(conn *amqp.Connection) ConnectionInfo {
	return ConnectionInfo{
		Vhost:      conn.Config.Vhost,
		Major:      conn.Major,
		Minor:      conn.Minor,
		Locales:    conn.Locales,
		Server:     conn.Properties,
		ChannelMax: conn.Config.ChannelMax,
		FrameSize:  conn.Config.FrameSize,
		Heartbeat:  conn.Config.Heartbeat,
		Locale:     conn.Config.Locale,
		Properties: conn.Config.Properties,
		TLS:        conn.Config.TLSClientConfig,
	}
}

func (c ConnectionInfo) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

// QueueInfo containes information about an AMQP Queue.
type QueueInfo struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Consumers int    `json:"consumers" yaml:"consumers"`
	Messages  int    `json:"messages" yaml:"messages"`
}

// GetInfoFromQueue extracts returns a QueueInfo struct from an AMQP Queue.
func GetInfoFromQueue(q amqp.Queue) QueueInfo {

	return QueueInfo{
		Name:      q.Name,
		Consumers: q.Consumers,
		Messages:  q.Messages,
	}
}

func (q QueueInfo) String() string {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
