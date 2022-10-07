package main


import (
    "fmt"
    "sync"
)

type Person struct {
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
            // read only: ensure immutable
            fmt.Printf("%v: %v\n", i, p)
        }(i)
    }
    go func() {
        wg.Wait()
        close(ch)
    }()
    return ch
}

func main() {
    var p = &Person{"Alice", "xxx", 0}
    done := run(5, p)
    <-done
    fmt.Printf("%v\n", p)
}
