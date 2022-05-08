package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
	"sync"
	"time"
)

const (
	RabbitMQDefaultHost        = "127.0.0.1"
	RabbitMQDefaultUsername    = "guest"
	RabbitMQDefaultPassword    = "guest"
	RabbitMQDefaultPort        = 5672
	RabbitMQDefaultVirtualHost = "/"
	RabbitMQDefaultRoutingKey  = "#"
)

type RabbitmqClient struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	cacheExchange []*Exchange
	cacheQueue    []*Queue
	cacheBinding  []*Binding
	cacheConsumer []*Consumer

	// true 连接成功、 false 重连中
	recoveryFlag bool
	recoveryLock *sync.RWMutex

	postFun []func()
}

const (
	Direct  = "direct"
	Topic   = "topic"
	Fanout  = "fanout"
	Headers = "headers"
)

func (rc *RabbitmqClient) Init() {
	rc.cacheExchange = make([]*Exchange, 0, 0)
	rc.cacheQueue = make([]*Queue, 0, 0)
	rc.cacheBinding = make([]*Binding, 0, 0)
	rc.cacheConsumer = make([]*Consumer, 0, 0)

	rc.recoveryLock = new(sync.RWMutex)
}

func (rc *RabbitmqClient) IsConnect() bool {
	return rc.recoveryFlag
}

/*
| ----------------------------------------------|
                 recovery connect
| ----------------------------------------------|
*/

func (rc *RabbitmqClient) Connect(rmd *RabbitMQMetadata) error {
	if rmd == nil {
		return errors.New("rabbitMQ connect metadata must not blank")
	}
	rmd.verityAndFillDefault()
	go rc.retryConnectRabbitMQ(rmd)
	return nil
}

func (rc *RabbitmqClient) retryConnectRabbitMQ(rmd *RabbitMQMetadata) {
	uri := rmd.buildConnURI()
	for range time.Tick(5 * time.Second) {
		logrus.WithFields(logrus.Fields{
			"do": fmt.Sprintf("try to connect rabbitmq server(%s:%d)", rmd.Host, rmd.Port),
		}).Info()

		if rc.connection == nil {
			// build tcp connect
			conn, err := amqp.Dial(uri)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"do":    fmt.Sprintf("build rabbitmq server(%s:%d) connection fail", rmd.Host, rmd.Port),
					"cause": fmt.Sprintf("%v", err),
				}).Error()
				continue
			}
			rc.connection = conn
			// recovery
			exceptionConn := conn.NotifyClose(make(chan *amqp.Error))
			go func() {
				for connErr := range exceptionConn {
					rc.connection = nil
					logrus.WithFields(logrus.Fields{
						"do":    fmt.Sprintf("receive rabbitmq server(%s:%d) connection disconnect notification", rmd.Host, rmd.Port),
						"cause": fmt.Sprintf("%v", connErr),
					}).Error()
					if rc.isNeedRecovery() {
						go rc.retryConnectRabbitMQ(rmd)
					}
				}
			}()
		}

		// build channel connect
		channel, err := rc.connection.Channel()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"do":    fmt.Sprintf("build rabbitmq server(%s:%d) channel fail", rmd.Host, rmd.Port),
				"cause": fmt.Sprintf("%v", err),
			}).Error()
			continue
		}
		// recovery
		exceptionChn := channel.NotifyClose(make(chan *amqp.Error))
		go func() {
			for connErr := range exceptionChn {
				rc.channel = nil
				logrus.WithFields(logrus.Fields{
					"do":    fmt.Sprintf("receive rabbitmq server(%s:%d) channel disconnect notification", rmd.Host, rmd.Port),
					"cause": fmt.Sprintf("%v", connErr),
				}).Error()
				if rc.isNeedRecovery() {
					go rc.retryConnectRabbitMQ(rmd)
				}
			}
		}()
		err = channel.Qos(250, 0, false)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"do":    "set channel qos fail",
				"cause": fmt.Sprintf("%v", err),
			}).Error()
			continue
		}

		rc.channel = channel
		rc.recoveryFlag = true

		// reload metadata
		rc.playbackCacheMetadata()

		logrus.WithFields(logrus.Fields{
			"do": fmt.Sprintf("connect rabbitmq server(%s:%d) success", rmd.Host, rmd.Port),
		}).Info()
		break
	}
}

func (rc *RabbitmqClient) isNeedRecovery() bool {
	rc.recoveryLock.Lock()
	defer rc.recoveryLock.Unlock()

	if rc.recoveryFlag {
		rc.recoveryFlag = false
		return true
	} else {
		return false
	}
}

func (rc *RabbitmqClient) playbackCacheMetadata() {
	if rc.cacheExchange != nil && len(rc.cacheExchange) > 0 {
		for _, exchange := range rc.cacheExchange {
			err := rc.DeclareExchange(exchange)
			if err != nil {
				fmt.Println(fmt.Sprintf("declare Exchange fail, cause: %v", err))
			}
		}
	}
	if rc.cacheQueue != nil && len(rc.cacheQueue) > 0 {
		for _, queue := range rc.cacheQueue {
			err := rc.DeclareQueue(queue)
			if err != nil {
				fmt.Println(fmt.Sprintf("declare Queue fail, cause: %v", err))
			}
		}
	}
	if rc.cacheBinding != nil && len(rc.cacheBinding) > 0 {
		for _, binding := range rc.cacheBinding {
			err := rc.DeclareBinding(binding)
			if err != nil {
				fmt.Println(fmt.Sprintf("declare consumer fail, cause: %v", err))
			}
		}
	}
	if rc.cacheConsumer != nil && len(rc.cacheConsumer) > 0 {
		for _, consumer := range rc.cacheConsumer {
			err := rc.RegisterConsumer(consumer)
			if err != nil {
				fmt.Println(fmt.Sprintf("register consumer fail, cause: %v", err))
			}
		}
	}
}

/*
| ----------------------------------------------|
                     declare
| ----------------------------------------------|
*/

func (rc *RabbitmqClient) DeclareExchange(exchange *Exchange) error {
	if exchange == nil || !exchange.verityAndFillDefault() {
		return errors.New("invalid parameter")
	}

	// cache
	if !exchange.cacheFlag {
		exchange.cacheFlag = true
		rc.cacheExchange = append(rc.cacheExchange, exchange)
	}

	if rc.channel == nil {
		return nil
	}

	return rc.channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args)
}

func (rc *RabbitmqClient) DeclareQueue(queue *Queue) error {
	if queue == nil || !queue.verityAndFillDefault() {
		return errors.New("invalid parameter")
	}

	// cache
	if !queue.cacheFlag {
		queue.cacheFlag = true
		rc.cacheQueue = append(rc.cacheQueue, queue)
	}

	if rc.channel == nil {
		return nil
	}

	_, err := rc.channel.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args)
	return err
}

func (rc *RabbitmqClient) DeclareBinding(binding *Binding) error {
	if binding == nil {
		return errors.New("invalid parameter")
	}

	// cache
	if !binding.cacheFlag {
		binding.cacheFlag = true
		rc.cacheBinding = append(rc.cacheBinding, binding)
	}

	if rc.channel == nil {
		return nil
	}

	return rc.channel.QueueBind(binding.Queue, binding.RoutingKey, binding.Exchange, binding.NoWait, binding.Args)
}

/*
| ----------------------------------------------|
                register consumer
| ----------------------------------------------|
*/

func (rc *RabbitmqClient) RegisterConsumer(consumer *Consumer) error {
	if len(strings.TrimSpace(consumer.Queue)) == 0 || consumer.ConsumeFun == nil {
		return errors.New("invalid parameter")
	}

	// cache
	if !consumer.cacheFlag {
		consumer.cacheFlag = true
		rc.cacheConsumer = append(rc.cacheConsumer, consumer)
	}

	if rc.channel == nil {
		return nil
	}

	deliveries, err := rc.channel.Consume(consumer.Queue, consumer.ConsumerTag, consumer.AutoAck, consumer.Exclusive, consumer.NoLocal, consumer.NoWait, consumer.Args)
	if err != nil {
		return err
	}
	go func() {
		for delivery := range deliveries {
			go dealRegisterFun(delivery, consumer.ConsumeFun, consumer.AutoAck)
		}
	}()
	return nil
}

func dealRegisterFun(delivery amqp.Delivery, consumeFun func(delivery *amqp.Delivery), autoAck bool) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"do":    fmt.Sprintf("deal rabbitmq message(%s) fail", string(delivery.Body)),
				"cause": fmt.Sprintf("%v", err),
			}).Error()
			if !autoAck {
				err := delivery.Nack(false, true)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"do":    fmt.Sprintf("nack rabbitmq message(%s) fail", string(delivery.Body)),
						"cause": fmt.Sprintf("%v", err),
					}).Error()
				}
			}
		} else {
			if !autoAck {
				err = delivery.Ack(false)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"do":    fmt.Sprintf("ack rabbitmq message(%s) fail", string(delivery.Body)),
						"cause": fmt.Sprintf("%v", err),
					}).Error()
				}
			}
		}
	}()
	consumeFun(&delivery)
}

/*
| ----------------------------------------------|
                    publish
| ----------------------------------------------|
*/

func (rc *RabbitmqClient) Publish2(exchange, key string, mandatory, immediate bool, publish amqp.Publishing) error {
	if rc.channel == nil {
		return errors.New("the connection task is not complete")
	}
	if err := rc.channel.Publish(exchange, key, mandatory, immediate, publish); err != nil {
		return err
	}
	return nil
}

func (rc *RabbitmqClient) PublishJson(exchange string, key string, obj interface{}) error {
	if obj == nil {
		return errors.New("invalid parameter")
	}
	if rc.channel == nil {
		return errors.New("the connection task is not complete")
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = rc.channel.Publish(exchange, key, false, false, amqp.Publishing{
		Timestamp:   time.Now(),
		ContentType: "application/json",
		Body:        jsonBytes,
	})
	return err
}
