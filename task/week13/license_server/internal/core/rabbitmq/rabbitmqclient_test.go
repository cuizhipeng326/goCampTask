package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func Test_RabbitMQConsumer(t *testing.T) {
	rabbitMqClient := new(RabbitmqClient)
	rabbitMqClient.Init()

	err := rabbitMqClient.Connect(&RabbitMQMetadata{
		Host:        "192.168.9.247",
		Port:        51887,
		Username:    "admin",
		Password:    "123456",
		VirtualHost: "/",
	})
	if err != nil {
		return
	}

	declareQueueAndBinding(rabbitMqClient, "exchange_test", "queue_test1", "key_test")
	declareQueueAndBinding(rabbitMqClient, "exchange_test", "queue_test2", "key_test")

	err = rabbitMqClient.RegisterConsumer(&Consumer{
		Queue:       "queue_test1",
		ConsumerTag: "测试队列1",
		ConsumeFun: func(delivery *amqp.Delivery) {
			fmt.Println("queue_test1 receive: ", string(delivery.Body))
			if err != nil {
				fmt.Println("queue_test1 ack fail, cause: ", fmt.Sprintf("%v", err))
			}
			time.Sleep(time.Second)
		},
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("register consumer fail, cause: %v", err))
	}

	err = rabbitMqClient.RegisterConsumer(&Consumer{
		Queue:       "queue_test2",
		ConsumerTag: "测试队列2",
		ConsumeFun: func(delivery *amqp.Delivery) {
			fmt.Println("queue_test2 receive: ", string(delivery.Body))
			time.Sleep(time.Second)
			panic("custom throw exception")
		},
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("######## register consumer fail, cause: %v", err))
	}

	select {}
}

func Test_RabbitMQPublish(t *testing.T) {
	rabbitMqClient := new(RabbitmqClient)
	rabbitMqClient.Init()

	err := rabbitMqClient.Connect(&RabbitMQMetadata{
		Host:        "192.168.9.247",
		Port:        51887,
		Username:    "admin",
		Password:    "123456",
		VirtualHost: "/",
	})
	if err != nil {
		return
	}

	err = rabbitMqClient.DeclareExchange(&Exchange{
		Name: "exchange_test",
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("######## declare exchange fail, cause: %v", err))
	}

	obj := &Obj{
		Id:   1,
		Name: "zy",
		Age:  25,
		Sex:  0,
		City: "甘肃",
	}

	//for now := range time.Tick(time.Millisecond) {
	//	obj.TimeStamp = now.Unix()
	//	err = rabbitMqClient.PublishJson("exchange_test", "key_test", obj)
	//	if err != nil {
	//		fmt.Println(fmt.Sprintf("######## publish message fail, cause: %v", err))
	//	}
	//}

	time.Sleep(10 * time.Second)

	obj.TimeStamp = time.Now().Unix()
	err = rabbitMqClient.PublishJson("exchange_test", "key_test", obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("######## publish message fail, cause: %v", err))
	}
}

type Obj struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Sex       int    `json:"sex"`
	City      string `json:"city"`
	TimeStamp int64  `json:"timeStamp"`
}

func declareQueueAndBinding(rabbitMqClient *RabbitmqClient, exchange, queue, routingKey string) {
	err := rabbitMqClient.DeclareQueue(&Queue{
		Name: queue,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("######## declare queue fail, cause: %v", err))
	}
	err = rabbitMqClient.DeclareBinding(&Binding{
		Exchange:   exchange,
		Queue:      queue,
		RoutingKey: routingKey,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("######## declare binding fail, cause: %v", err))
	}
}
