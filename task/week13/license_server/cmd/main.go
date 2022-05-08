package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/go-kratos/kratos/pkg/naming/nacos"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver"
	"license_kratos/internal/core/args"
	"license_kratos/internal/di"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

//go:generate go build -o ../target/centos/license_server/license_server ./cmd
const AppId = "license-server"

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("license_server start")
	err := paladin.Init()
	if err != nil {
		panic(err)
	}
	// init log
	var cfg struct {
		Log   *log.Config
		Nacos *nacos.Config
	}
	if err = paladin.Get("log.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	log.Init(cfg.Log) // debug flag: log.dir={path}
	defer log.Close()
	// resolver
	if err = paladin.Get("nacos.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}

	if args.RunArgs.MasterIp != "" {
		for i := range cfg.Nacos.ServerConfigs {
			cfg.Nacos.ServerConfigs[i].IpAddr = args.RunArgs.MasterIp
		}
	}
	if err != nil {
		panic(err)
	}

	// http grpc 配置
	var (
		grpcCfg  warden.ServerConfig
		httpCfg  warden.ServerConfig
		ct       paladin.TOML
		grpcAddr string
		httpAddr string
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&grpcCfg); err != nil {
		return
	}
	grpcAddr = fmt.Sprintf("grpc://%s", grpcCfg.Addr)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&httpCfg); err != nil {
		return
	}
	httpAddr = fmt.Sprintf("http://%s", httpCfg.Addr)

	if args.RunArgs.LocalIp != "" {
		addrReg := regexp.MustCompile("(http|grpc):\\/\\/.*:(.*)")
		s := fmt.Sprintf("${1}://%s:${2}", args.RunArgs.LocalIp)
		grpcAddr = addrReg.ReplaceAllString(grpcAddr, s)
		httpAddr = addrReg.ReplaceAllString(httpAddr, s)
	}

	resolver.Register(nacos.Builder(cfg.Nacos))
	// 启动程序
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}
	//注册服务
	hn, _ := os.Hostname()
	instance := &naming.Instance{
		Region:   env.Region,
		Zone:     env.Zone,
		Env:      env.AppID,
		AppID:    AppId,
		Hostname: hn,
		Addrs:    []string{httpAddr},
		Version:  "v3.0",
		Metadata: map[string]string{"weight": "10"},
		Status:   1,
	}
	cancelFunc, err := nacos.Register(nil, context.Background(), instance)
	if err != nil {
		panic(err)
	}
	defer cancelFunc()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			err := nacos.Builder(nil).Close()
			if err != nil {
				panic(err)
			}
			log.Info("license_server exit")
			time.Sleep(time.Second)
			os.Exit(0)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
