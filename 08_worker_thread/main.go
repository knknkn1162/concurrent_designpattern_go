package main

import (
    "fmt"
    "time"
    "math/rand"
    "sync"
)

type Data struct {
    Number int
    CreatedBy string
}

func do_calc(num int) {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

func request(lists... string) <-chan Data{
    res := make(chan Data)
    var wg sync.WaitGroup
    for idx, str := range lists {
        wg.Add(1)
        go func(i int, str string) {
            defer wg.Done()
            for {
                d := Data{i, str}
                fmt.Printf("request -> %v\n", d)
                res <- Data{i, str}
                do_calc(1000)
            }
        }(idx, str)
    }
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func worker(ch <-chan Data, jobNum int) <-chan bool{
    done := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for val := range ch {
                fmt.Printf("worker #%v <- %v\n", i, val)
                do_calc(10000)
            }
        }(i)
    }
    go func() {
        defer close(done)
        wg.Wait()
    }()
    return done
}

func main() {
    ch := request("Alice", "Bobby", "Chris")
    done := worker(ch, 100)
    <-done
}
