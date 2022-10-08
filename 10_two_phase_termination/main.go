package main

import (
    "fmt"
    "math/rand"
    "time"
)

func do_calc(num int) int {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
    return m
}

// 1st phase: inform users of time-out
func monitor(num int) <-chan bool {
    done := make(chan bool)
    timeout := time.After(time.Duration(num) * time.Second)
    go func() {
        defer close(done)
        <-timeout
        fmt.Println("[1st] timeout!")
    }()
    return done
}

func countUp(num int) <-chan bool {
    res := make(chan bool)
    // monitor timeout
    done := monitor(num)
    count := 0
    go func() {
        defer close(res)
LOOP:
        for {
            ans :=  do_calc(1000)
            count++
            fmt.Printf("continue countup %02d (time: %v[ms])\n", count, ans)
            select {
            case <-done:
                fmt.Println("[2nd] recv timeout -> terminate")
                break LOOP
            default:
            }
        }
    }()
    return res
}

func main() {
    <-countUp(10)
    fmt.Println("finished!")
}
