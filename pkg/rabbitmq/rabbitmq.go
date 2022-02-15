package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/setarek/pym-particle-microservice/config"
	"github.com/setarek/pym-particle-microservice/pkg/logger"

	"github.com/streadway/amqp"
)

const (
	ShortenerQueue = "redirect"
	Visitors = "visitors"
)

type rabbitMQ struct {
	protocol    string
	user        string
	password    string
	host        string
	port        string
	conn        *amqp.Connection

}

var Channel *amqp.Channel


func InitRabbitMQ(config *config.Config, logger logger.Logger) {

	conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%v/", config.GetString("rabbitmq_protocol"),
		config.Get("rabbitmq_user"), config.GetString("rabbitmq_password"), config.GetString("rabbitmq_host"),
		config.Get("rabbitmq_port")))

	// todo: implement connection recovery for channel and consumers
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("error while openning a channel", err)
	}

	_, err = ch.QueueDeclare(ShortenerQueue, true, false, false, false, nil)
	if err != nil {
		logger.Error("error while declaring a queue", err)
	}

	_, err = ch.QueueDeclare(Visitors, true, false, false, false, nil)
	if err != nil {
		logger.Error("error while declaring a queue", err)
	}
	Channel = ch
}

func GetQueue()  (channel *amqp.Channel){
	return Channel
}

func PublishDurableMessage(ctx context.Context, logger logger.Logger,queueName string, message map[string]interface{})  {
	fmt.Println("Ssssssssssssssssssssssssssssssssssssssssssssssssssssss")
	span, _ := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("rabbitmq.PublishDurableMessage.%s", queueName))
	defer span.Finish()
	rabbitMQClient := GetQueue()
	marshalled, err := json.Marshal(message)
	err = rabbitMQClient.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         marshalled,
		})
	if err != nil {
		logger.Error("error while publish message", err)
	}
}
