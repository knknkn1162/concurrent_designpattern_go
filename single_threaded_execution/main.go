package main


import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Person struct {
    Name string
    History string
    Count  int32
}

func run(jobNum int, p *Person) <-chan bool{
    ch := make(chan bool)
    var wg sync.WaitGroup
    var mx sync.Mutex
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func() {
            wg.Done()
            // single threaded execution with mutex
            mx.Lock()
            p.History += "x"
            mx.Unlock()
            // single threaded execution with atomic
            atomic.AddInt32(&p.Count, 1)
        }()
    }
    go func() {
        wg.Wait()
        close(ch)
    }()
    return ch
}

func main() {
    var p = &Person{"Alice", "", 0}
    done := run(5, p)
    <-done
    fmt.Printf("%v\n", p)
}
