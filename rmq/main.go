package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/streadway/amqp"
)

func main() {
	port := os.Getenv("PORT")
	serviceName := os.Getenv("RMQ_SERVICENAME")
	handler := WebHandler(NewRMQClient, serviceName)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

type RMQClient interface {
	Connect(uri, channelName string) error
	Send(value string) error
	Receive() (string, error)
	Close()
}

type RMQClientFactory func() RMQClient

func WebHandler(fac RMQClientFactory, serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := PerformTest(fac, serviceName); err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			fmt.Fprintf(w, "Failed to access RMQ: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "RMQ service is OK")
	}
}

func PerformTest(fac RMQClientFactory, serviceName string) error {
	client := fac()
	defer client.Close()
	_, uri, err := GetURI(serviceName)
	if err != nil {
		return err
	}
	channelName := "aChannel"
	value := "a value"
	if err := client.Connect(uri, channelName); err != nil {
		return err
	}
	if err := client.Send(value); err != nil {
		return err
	}
	if newValue, err := client.Receive(); err != nil {
		return err
	} else if newValue != value {
		return errors.New("Did not receive back the value I sent")
	}
	return nil
}

func GetURI(serviceName string) (ssl bool, uri string, err error) {
	app, err := cfenv.Current()
	if err != nil {
		return
	}

	svc, err := app.Services.WithName(serviceName)
	if err != nil {
		return
	}

	ssl, _ = svc.Credentials["ssl"].(bool)
	uri, _ = svc.CredentialString("uri")

	return
}

type RMQClientImpl struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewRMQClient() RMQClient {
	return &RMQClientImpl{}
}

func (c *RMQClientImpl) Connect(uri string, channelName string) (err error) {
	c.conn, err = amqp.Dial(uri)
	if err != nil {
		return
	}
	c.ch, err = c.conn.Channel()
	if err != nil {
		return
	}
	c.q, err = c.ch.QueueDeclare(channelName, false, false, false, false, nil)
	if err != nil {
		return
	}
	return
}

func (c *RMQClientImpl) Send(value string) error {
	return c.ch.Publish(
		"",
		c.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(value),
		},
	)
}

func (c *RMQClientImpl) Receive() (value string, err error) {
	msgs, err := c.ch.Consume(c.q.Name, "", true, false, false, false, nil)
	if err != nil {
		return
	}
	msg := <-msgs
	value = string(msg.Body)
	return
}

func (c *RMQClientImpl) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	if c.ch != nil {
		c.ch.Close()
	}
}
