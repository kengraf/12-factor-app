package main

import (
    "fmt"
    "net/http"
    "math/rand"
    "time"
)

func hello(w http.ResponseWriter, req *http.Request) {

    time.Sleep( rand.Intn(10) * time.Second)
    fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func dowork(w http.ResponseWriter, req *http.Request) {

    var sum = 0
    for i := 1; i < 50000; i++ {
        for j := 1; j < i; j++ {
            sum += i * rand.Intn(j)
        }
    }
    fmt.Fprintf(w, "done %d\n", sum )
}

func main() {

    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/dowork", dowork)

    http.ListenAndServe(":8090", nil)
}

