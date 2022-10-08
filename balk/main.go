package main

import (
    "fmt"
    "os"
    "math/rand"
    "time"
)

type FileInfo struct {
    Fname string
    Content string
    Changed bool
}

func (finfo *FileInfo) save(content string) {
    file, err := os.OpenFile(finfo.Fname, os.O_CREATE|os.O_WRONLY, 0644)
    defer file.Close()
    if err != nil {
        panic(err)
    }
    if content != "" {
        finfo.update(content)
    } else {
        if finfo.Changed {
            file.WriteString(finfo.Content)
            fmt.Printf("saved!\n")
        } else {
            fmt.Printf("save skipped!\n")
        }
    }
    finfo.Changed = false
}

func (finfo *FileInfo) update(content string) {
    finfo.Content = content
    finfo.Changed = true
}

func initialize(fname string) *FileInfo {
    os.Remove(fname)
    finfo := &FileInfo{fname, "", false}
    finfo.save(fmt.Sprintf("start: %03d", 0))
    fmt.Println("initialize!")
    return finfo
}

func saver(finfo *FileInfo, p time.Duration) <-chan bool{
    res := make(chan bool)
    go func() {
        defer close(res)
        ticker := time.Tick(p)
        for {
            fmt.Printf("saver: %v (%v)\n", finfo.Changed, finfo.Content)
            finfo.save("")
            <-ticker
        }
    }()
    return res
}

func changer(finfo *FileInfo) <-chan bool{
    res := make(chan bool)
    go func() {
        defer close(res)
        for i := 1; ; i++ {
            // might change
            fmt.Printf("changer: %v (%v)\n", finfo.Changed, finfo.Content)
            finfo.update(fmt.Sprintf("Start: %03d", i))
            m := rand.Intn(1000) + 500
            time.Sleep(time.Duration(m) * time.Millisecond)
        }
    }()
    return res
}

func main() {
    finfo := initialize("data.txt")
    saver(finfo, 3 * time.Second)
    ch := changer(finfo)
    <-ch
    fmt.Println("done")
}
