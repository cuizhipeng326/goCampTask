module license_kratos

go 1.13

require (
	github.com/go-kratos/kratos v0.6.1-0.20210110073301-6e49fe9ac61e
	github.com/go-redis/redis/v8 v8.4.2 // indirect
	github.com/go-resty/resty/v2 v2.4.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/wire v0.5.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.5.0
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/wenzhenxi/gorsa v0.0.0-20191231021121-58a13482fb09
	google.golang.org/genproto v0.0.0-20200402124713-8ff61da6d932
	google.golang.org/grpc v1.29.1
)

replace github.com/go-kratos/kratos => ../../golang_store/src/github.com/go-kratos/kratos@v0.6.0
