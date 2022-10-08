package main

import (
    "fmt"
    "math/rand"
    "time"
    "sync"
)
type Table struct {
    CreatedBy string
    Number int
}

func do_sleep(num int) {
    m := rand.Intn(num)
    time.Sleep(time.Duration(m) * time.Millisecond)
}

func maker(num int) <-chan Table{
    res := make(chan Table, num)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for j := 0; ; j++ {
                do_sleep(500)
                tbl := Table{
                    fmt.Sprintf("MakerThread-%d", i),
                    j,
                }
                fmt.Printf("maker #%03d: ->(#%03d, createdBy: %v)\n", i, tbl.Number, tbl.CreatedBy)
                res <- tbl
            }
        }(i)
    }
    go func() {
        wg.Wait()
        close(res)
    }()
    return res
}

func eater(num int, ch <-chan Table) <-chan bool {
    res := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < num; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for val := range ch {
                do_sleep(500)
                fmt.Printf("eater #%03d: <-(#%03d, createdBy: %v)\n", i, val.Number, val.CreatedBy)
            }
        }(i)
    }
    go func() {
        wg.Wait()
        close(res)
    }()
    return res

}

func main() {
    tblCh := maker(3)
    done := eater(3, tblCh)
    <-done
}
