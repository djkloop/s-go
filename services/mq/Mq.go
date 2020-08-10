package mq

import (
	"bytes"
	"fmt"

	"github.com/streadway/amqp"
)

type CallBack func(msg string)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	return conn, err
}

// 发送端函数
func Publish(exchangeName string, queueName string, body string) error {
	//
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建一个通道
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 创建队列
	q, err := channel.QueueDeclare(
		queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// 发送消息
	err = channel.Publish(exchangeName, q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}

// 接收者方法
func Consumer(exchangeName string, queueName string, callback CallBack) {
	// 建立连接
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建通道
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取数据
	messages, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			s := BytesToString(&(d.Body))
			callback(*s)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for messages")
	<-forever
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

func PublishEx(exchangeName string, types string, routingKey string, body string) error {
	// 建立连接
	conn, err := Connect()
	defer conn.Close()

	if err != nil {
		return err
	}

	// 创建channel
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(exchangeName, types, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}

func ConsumerEx(exchangeName string, types string, routingKey string, callback CallBack) {
	// 建立连接
	conn, err := Connect()
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	//创建通道
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建交换机
	err = channel.ExchangeDeclare(exchangeName, types, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个队列
	q, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 绑定
	err = channel.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 接收消息
	messages, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			s := BytesToString(&(d.Body))
			callback(*s)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for message\n")
	<-forever
}

// 死信队列消费端
func ConsumerDlx(exchangeA string, queueAName string, exchangeB string, queueBName string, ttl int, callback CallBack) {
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建channel
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建A交换机
	// A队列
	// A交换机和A队列
	err = channel.ExchangeDeclare(exchangeA, "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	queueA, err := channel.QueueDeclare(
		queueAName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":          ttl,
			"x-dead-letter-exchange": exchangeB,
			//"x-dead-letter-queue": "",
			//"x-dead-letter-routeing-key": "key"
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.QueueBind(queueA.Name, "", exchangeA, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建B交换机
	// 创建B队列
	// 绑定b交换机和队列
	err = channel.ExchangeDeclare(exchangeB, "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	queueB, err := channel.QueueDeclare(
		queueBName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.QueueBind(queueB.Name, "", exchangeB, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	messages, err := channel.Consume(queueB.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			s := BytesToString(&(d.Body))
			callback(*s)
			_ = d.Ack(false)
		}
	}()
	fmt.Println("Waiting for messages")
	<-forever
}

func PublishDlx(exchangeA string, body string) error {
	// 建立连接
	conn, err := Connect()
	defer conn.Close()

	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil {
		return err
	}

	err = channel.Publish(exchangeA, "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}
