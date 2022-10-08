package main

import (
    "fmt"
    "time"
    "sync"
    "math/rand"
)

type Person struct {
    mx sync.RWMutex
    Name string
    History string
    Count int32
}

func do_calc(num int) {
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

func reader(p *Person, num int) <-chan bool{
    res := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for {
                // reader lock
                p.mx.RLock()
                str := p.History
                p.mx.RUnlock()

                fmt.Printf("reader #%v: %v\n", i, str)
                do_calc(1000)
            }
        }(i)
    }
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func writer(p *Person, num int) <-chan bool{
    res := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for {
                // write lock
                p.mx.Lock()
                p.Count += 1
                str := fmt.Sprintf("%05d", p.Count)
                p.History = str
                p.mx.Unlock()

                fmt.Printf("writer #%v: %v\n", i, str)
                do_calc(1000)
            }
        }(i)
    }
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func main() {
    p := &Person{sync.RWMutex{}, "test", "*", 0}
    reader(p, 6)
    ch := writer(p, 2)
    <-ch
}
