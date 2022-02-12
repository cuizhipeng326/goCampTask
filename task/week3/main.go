// 作业题：基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
package main

import (
    "context"
    "fmt"
    "golang.org/x/sync/errgroup"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    sigChan := make(chan os.Signal, 1) // 带缓冲是为了避免收到停止信号主动close掉server后再触发server的协程通过sigChan通知监控信号协程时不阻塞
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

    ctx := context.Background()

    group := errgroup.Group{}
    group.Go(func() error {
        _, cancel := context.WithCancel(ctx)
        defer func() {
            cancel() // 利用context实现通知监控信号的协程退出
            // sigChan <- syscall.SIGTERM // 复用带缓冲的sigChan通知监控停止信号的协程退出，不够优雅
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
