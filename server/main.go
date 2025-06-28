package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/bytemare/opaque"
)

func main() {
    http.HandleFunc("/auth", handleAuth)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
    var msg struct {
        Protocol string `json:"protocol"`
        Step     string `json:"step"`
        Data     string `json:"data"`
    }
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        http.Error(w, "Invalid request", 400)
        return
    }

    if msg.Protocol == "opaque" && msg.Step == "init" {
        server := opaque.NewServer()
        response, err := server.DeserializeRegistrationRequest([]byte(msg.Data))
        if err != nil {
            http.Error(w, "Failed to process", 500)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(struct {
            Data string `json:"data"`
        }{Data: string(response)})
    }
}
