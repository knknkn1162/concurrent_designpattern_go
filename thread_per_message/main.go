package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Data struct {
    Message string
    Number int
}

func do_calc(num int) {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

func request(num int) <-chan Data {
    res := make(chan Data)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            do_calc(1000)
            fmt.Printf("request %v %v\n", string(i+0x41), i)
        }(i)
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
