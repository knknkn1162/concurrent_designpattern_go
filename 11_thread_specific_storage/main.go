package main

import (
    "fmt"
    "time"
    "log"
    "os"
    "io"
    "sync"
    "math/rand"
)

func do_calc(num int) int {
    rand.Seed(time.Now().UnixNano())
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
    return m
}

func logger() <-chan bool{
    done := make(chan bool)
    var wg sync.WaitGroup
    lists := []string{"Alice", "Bobby", "Chris"}
    for _, val := range lists {
        wg.Add(1)
        go func(id string) {
            file, err := os.OpenFile(id + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
            if err != nil {
                panic(err)
            }
            defer file.Close()
            defer wg.Done()
            multiLogFile := io.MultiWriter(os.Stdout, file)
            log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
            log.SetOutput(multiLogFile)
            // log for each file
            do_calc(300)
            log.Printf("%v says...\n", id)
        }(val)
    }
    go func() {
        defer close(done)
        wg.Wait()
    }()
    return done
}

func main() {
    <-logger()
    fmt.Println("done!")
}
