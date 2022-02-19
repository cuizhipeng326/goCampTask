// +build wireinject

package main

import (
    . "demo1/internal/pkg/job"
    "github.com/google/wire"
)

func SuperFactory() *Factory {
    panic(wire.Build(NewFactory, NewBuilder, NewWorker, NewWorkFunc))
    // return &Factory{}
}
