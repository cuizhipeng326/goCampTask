package main

import (
    "fmt"
    "net/http"
)

var (
    httpServer *http.Server
)

func init() {

    httpServer = new(http.Server)
    httpServer.Addr = "127.0.0.1:12345"
}

func Serve() error {
    err := httpServer.ListenAndServe()
    fmt.Printf("%+v\n", err)
    return err
}

func Close() error {
    return httpServer.Close()
}
