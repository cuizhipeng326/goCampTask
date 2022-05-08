package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"net/url"
	"strings"
	"time"
)

type RabbitMQMetadata struct {
	Host        string
	Port        int
	Username    string
	Password    string
	VirtualHost string
}

func (rmd *RabbitMQMetadata) verityAndFillDefault() {
	if len(rmd.Host) == 0 {
		rmd.Host = RabbitMQDefaultHost
	}
	if rmd.Port == 0 {
		rmd.Port = RabbitMQDefaultPort
	}
	if len(rmd.Username) == 0 {
		rmd.Username = RabbitMQDefaultUsername
	}
	if len(rmd.Password) == 0 {
		rmd.Username = RabbitMQDefaultPassword
	}
}

func (rmd *RabbitMQMetadata) buildConnURI() string {
	rui := fmt.Sprintf("amqp://%s:%s@%s:%d/", url.QueryEscape(rmd.Username), url.QueryEscape(rmd.Password), rmd.Host, rmd.Port)
	if rmd.VirtualHost != "/" && len(rmd.VirtualHost) > 0 {
		rui = rui + rmd.VirtualHost
	}
	return rui
}

type Exchange struct {
	Type       string
	Name       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table

	cacheFlag bool
}

func (e *Exchange) verityAndFillDefault() bool {
	if len(strings.TrimSpace(e.Name)) == 0 {
		return false
	}
	if len(strings.TrimSpace(e.Type)) == 0 {
		e.Type = Topic
	}
	if e.Args == nil {
		e.Args = make(map[string]interface{})
		e.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
	} else {
		if _, exist := e.Args["x-message-ttl"]; !exist {
			e.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
		}
	}
	return true
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table

	cacheFlag bool
}

func (q *Queue) verityAndFillDefault() bool {
	if len(strings.TrimSpace(q.Name)) == 0 {
		return false
	}
	if q.Args == nil {
		q.Args = make(map[string]interface{})
		q.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
	} else {
		if _, exist := q.Args["x-message-ttl"]; !exist {
			q.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
		}
	}
	return true
}

type Binding struct {
	Exchange   string
	Queue      string
	RoutingKey string
	NoWait     bool
	Args       amqp.Table

	cacheFlag bool
}

func (b *Binding) ToString() string {
	return fmt.Sprintf("%s_%s_%s", b.Exchange, b.Queue, b.RoutingKey)
}

func (b *Binding) verityAndFillDefault() bool {
	if len(strings.TrimSpace(b.Exchange)) == 0 {
		return false
	}
	if len(strings.TrimSpace(b.Queue)) == 0 {
		return false
	}
	if len(strings.TrimSpace(b.RoutingKey)) == 0 {
		return false
	}
	if b.Args == nil {
		b.Args = make(map[string]interface{})
		b.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
	} else {
		if _, exist := b.Args["x-message-ttl"]; !exist {
			b.Args["x-message-ttl"] = 30 * time.Minute.Milliseconds()
		}
	}
	return true
}

type Consumer struct {
	Queue       string
	ConsumerTag string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	Args        amqp.Table
	ConsumeFun  func(delivery *amqp.Delivery)

	cacheFlag bool
}
