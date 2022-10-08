package main


import (
    "fmt"
    "time"
    "math/rand"
)

func client_run(num int) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        for i := 0; i < num; i++ {
            m := rand.Intn(300)
            time.Sleep(time.Duration(m) * time.Millisecond)
            str := fmt.Sprintf("%04d", i)
            fmt.Printf("client: %v\n", str)
            ch <- str
        }
    }()
    return ch
}

func server_run(ch <-chan string) <-chan bool{
    res := make(chan bool)
    go func() {
        defer close(res)
        // guard until the messsages are send via channel
        for str := range ch {
            m := rand.Intn(300)
            time.Sleep(time.Duration(m) * time.Millisecond)
            fmt.Printf("server: %v\n", str)
        }
    }()
    return res
}

func main() {
    // this is RequestQueue
    ch := client_run(100)
    done := server_run(ch)
    <-done
}
