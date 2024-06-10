package main

import (
    "log"
    "net/http"

)

const (
    port = "localhost:5005"
)

func main() {
    log.Println("[INFO]", "Connecting to db at", dbPort)
    d, err := Connect()
    if err != nil {
        log.Fatalln(err)
    }
    defer Disconnect(d)
    log.Println("[INFO]", "Connected to the database")

    mux := http.NewServeMux()

    sv := NewServer(port, d, mux)
    err = sv.DB.Init()
    if err != nil { log.Fatalln("[ERROR]", err) }

    Register(&sv)

    log.Println("[INFO]", "Starting server at port", sv.Addr)
    if err := http.ListenAndServe(sv.Addr, sv.Mux); err != nil {
        log.Fatalln("[ERROR]", err)
    }
}


