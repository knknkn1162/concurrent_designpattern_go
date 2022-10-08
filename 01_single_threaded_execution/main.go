package main


import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Person struct {
    mx sync.Mutex
    Name string
    History string
    Count  int32
}

func run(jobNum int, p *Person) <-chan bool{
    ch := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            fmt.Printf("run #%v\n", i)
            // single threaded execution with mutex
            p.mx.Lock()
            p.History += "x"
            p.mx.Unlock()
            // single threaded execution with atomic
            atomic.AddInt32(&p.Count, 1)
        }(i)
    }
    go func() {
        wg.Wait()
        close(ch)
    }()
    return ch
}

func main() {
    var p = &Person{sync.Mutex{}, "Alice", "", 0}
    done := run(5, p)
    <-done
    fmt.Printf("%v\n", p)
}
