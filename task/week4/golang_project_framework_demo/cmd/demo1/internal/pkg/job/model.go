package job

import "fmt"

type Work func() error

type Worker struct {
    Work
}

type Builder struct {
    worker *Worker
}

type Factory struct {
    builder *Builder
}

func (f *Factory) Run() {
    f.builder.worker.Work()
}

func NewWorkFunc() Work {
    return func() error {
        fmt.Println("start to work...")

        // do something...

    TAG:
        goto TAG

        return nil
    }
}

func NewWorker(workFunc Work) *Worker {
    return &Worker{workFunc}
}

func NewBuilder(worker *Worker) *Builder {
    return &Builder{worker}
}

func NewFactory(builder *Builder) *Factory {
    return &Factory{builder}
}
