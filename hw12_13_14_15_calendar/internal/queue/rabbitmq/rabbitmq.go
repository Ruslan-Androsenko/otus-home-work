package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeConf struct {
	Name        string `toml:"exchange_name"`
	Type        string `toml:"exchange_type"`
	Key         string `toml:"routing_key"`
	QueueName   string `toml:"queue_name"`
	ConsumerTag string `toml:"consumer_tag"`
}

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	dataSource string
	exchange   ExchangeConf
	deliveries <-chan amqp.Delivery
}

func New(dataSource string, exchange ExchangeConf) *RabbitMQ {
	return &RabbitMQ{
		dataSource: dataSource,
		exchange:   exchange,
	}
}

// Connect Подключение к очереди.
func (rmq *RabbitMQ) Connect() error {
	connection, err := amqp.Dial(rmq.dataSource)
	if err != nil {
		return fmt.Errorf("cannot Dial to Rabbit. Error: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("cannot open Channel on Rabbit. Error: %w", err)
	}

	if err = channel.ExchangeDeclare(
		rmq.exchange.Name,
		rmq.exchange.Type,
		true,
		false,
		false,
		true,
		nil,
	); err != nil {
		return fmt.Errorf("cannot declare Exchange on Rabbit. Error: %w", err)
	}

	if _, err = channel.QueueDeclare(
		rmq.exchange.QueueName,
		true,
		false,
		false,
		true,
		nil,
	); err != nil {
		return fmt.Errorf("cannot declare Queue on Rabbit. Error: %w", err)
	}

	if err = channel.QueueBind(
		rmq.exchange.QueueName,
		rmq.exchange.Key,
		rmq.exchange.Name,
		true,
		nil,
	); err != nil {
		return fmt.Errorf("cannot bind Channel on Queue. Error: %w", err)
	}

	rmq.connection = connection
	rmq.channel = channel

	return nil
}

func (rmq *RabbitMQ) Close() error {
	return rmq.connection.Close()
}

func (rmq *RabbitMQ) NotifyClose() chan *amqp.Error {
	return rmq.connection.NotifyClose(make(chan *amqp.Error))
}

// Publish Опубликовать сообщение в очередь.
func (rmq *RabbitMQ) Publish(message storage.Notification) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot Marshal message: %v. Error: %w", message, err)
	}

	if err = rmq.channel.Publish(
		rmq.exchange.Name,
		rmq.exchange.Key,
		false,
		false,
		amqp.Publishing{
			Body:            body,
			ContentType:     "application/json",
			ContentEncoding: "",
			Headers:         amqp.Table{},
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("cannot Publish to queue. Error: %w", err)
	}

	return nil
}

// Consume Получить сообщения из очереди.
func (rmq *RabbitMQ) Consume() error {
	deliveries, err := rmq.channel.Consume(
		rmq.exchange.QueueName,
		rmq.exchange.ConsumerTag,
		false,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		return fmt.Errorf("cannot Consume from queue. Error: %w", err)
	}

	rmq.deliveries = deliveries

	return nil
}

// Handle Обработка сообщений из очереди.
func (rmq *RabbitMQ) Handle(done chan error, logg *logger.Logger) {
	for d := range rmq.deliveries {
		logg.Infof(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		if err := d.Ack(false); err != nil {
			done <- fmt.Errorf("error Ack: %w", err)
		}
	}

	logg.Infof("Handle: deliveries channel closed")
	done <- nil
}
