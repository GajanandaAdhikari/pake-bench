package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"

    "github.com/bytemare/opaque"
)

func main() {
    client := opaque.NewClient()
    initMsg, err := client.SerializeRegistrationRequest()
    if err != nil {
        log.Fatal(err)
    }

    msg := struct {
        Protocol string `json:"protocol"`
        Step     string `json:"step"`
        Data     string `json:"data"`
    }{Protocol: "opaque", Step: "init", Data: string(initMsg)}

    body, _ := json.Marshal(msg)
    resp, err := http.Post("http://server:8080/auth", "application/json", bytes.NewBuffer(body))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    // Process server response...
}
