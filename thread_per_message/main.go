package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)


func do_calc(num int) {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

func request(num int) <-chan bool {
    res := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        fmt.Printf("req #%v begin\n", i)
        go func(i int) {
            defer wg.Done()
            // Assume that return value is not necessary
            fmt.Printf("goroutine #%v begin\n", i)
            do_calc(1000)
            fmt.Printf("goroutine #%v -> %v\n", string(i+0x41), i)
        }(i)
        fmt.Printf("req #%v end\n", i)
    }
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func main() {
    ch := request(3)
    <-ch
}
