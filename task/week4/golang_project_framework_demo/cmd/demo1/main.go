// cmd目录下的代码主要用于管理服务的生命周期
// google wire依赖注入实现服务配置等的创建
package main

import (
    "context"
    . "demo1/internal/pkg/server"
    "fmt"
    "golang.org/x/sync/errgroup"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    factory := SuperFactory()
    factory.Run()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

    group, ctx := errgroup.WithContext(context.Background())
    group.Go(func() error {
        _, cancel := context.WithCancel(ctx)
        defer func() {
            cancel() // 利用context实现通知监控信号的协程退出
            fmt.Println("server serve done")
        }()

        return Serve()
    })

    group.Go(func() error {
        select {
        case sig := <-sigChan: // 监控到停止信号，关闭server
            fmt.Printf("receive signal[%v], need to shut down server!\n", sig)
            return Close()
        case <-ctx.Done(): // server异常停止，不再等待停止信号，退出协程
            fmt.Println("server done!")
            return nil
        }
    })

    group.Wait()
}
