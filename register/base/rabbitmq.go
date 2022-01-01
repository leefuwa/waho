package base

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"waho/register"
)

type RabbitMqRegister struct {
	register.BaseRegister
}

var amq *amqp.Connection

func (rabbitMqRegister *RabbitMqRegister) Init() {
	isCheckAndNewRmq()
}

func Publish(queueName string, body string, exchangeDef ...string) error {
	isCheckAndNewRmq()

	channel, err := amq.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var exchange string
	if len(exchangeDef) > 0 {
		exchange = exchangeDef[0]
	}

	err = channel.Publish(exchange, queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}

func Consumer(queueName string, callback func(msg string,d *amqp.Delivery)) error {
	isCheckAndNewRmq()

	channel, err := amq.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()


	//创建queue
	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	deli, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range deli {
			s := bytes.NewBuffer(d.Body).String()
			d := d
			go func() {
				callback(s, &d)
			}()
		}
	}()
	log.Info("Wait ...")
	<-forever
	return err
}

func isCheckAndNewRmq()  {
	if amq == nil {
		amq, _ = NewRmq()
	}

	var err error
	_, err = amq.Channel()
	for err != nil {
		amq, _ = NewRmq()
		_, err = amq.Channel()
		// TODO 上报RMQ连接异常
	}
}

func NewRmq() (*amqp.Connection, error) {
	host := register.GetConf.Section("rabbitmq").GetDefault("host", "127.0.0.1")
	port := register.GetConf.Section("rabbitmq").GetDefault("port", "3306")
	user := register.GetConf.Section("rabbitmq").GetDefault("user", "root")
	pwd := register.GetConf.Section("rabbitmq").GetDefault("pwd", "")
	vhost := register.GetConf.Section("rabbitmq").GetDefault("vhost", "")

	url := "amqp://" + user + ":" + pwd + "@" + host + ":" + port + "/" + vhost
	amq, err := amqp.Dial(url)
	if err != nil {
		log.Panic("Rabbitmq Conn Error [" + url + "] :" + err.Error())
		return nil, err
	}
	return amq, err
}


