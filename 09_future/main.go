package main

import (
    "fmt"
    "strings"
    "math/rand"
    "sync"
    "time"
)


func do_calc(num int) {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

type Data struct {
    Number int
    Message string
}

func request(num int) <-chan Data {
    res := make(chan Data)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        fmt.Printf("req #%v begin\n", i)
        go func(i int) {
            defer wg.Done()
            fmt.Printf("req:goroutine #%v begin\n", i)
            do_calc(1000)
            // return value
            res <- Data{i, strings.Repeat(string(0x41+i), 10)}
            fmt.Printf("req:goroutine #%v\n", i)
        }(i)
        fmt.Printf("req #%v end\n", i)
    }
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func response(ch <-chan Data) <-chan bool {
    done := make(chan bool)
    go func() {
        defer close(done)
        // consume completed tasks in order
        for val := range ch {
            fmt.Printf("resp: %v\n", val)
        }
    }()
    return done
}

func main() {
    ch := request(3)
    <-response(ch)
}
