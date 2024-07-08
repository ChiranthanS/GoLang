package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "time"
)

type Response struct {
    ProcessedFile string `json:"processed_file"`
    Timestamp     string `json:"timestamp"`
}

func processFile(w http.ResponseWriter, r *http.Request) {
    filePath := r.URL.Query().Get("file")
    if filePath == "" {
        http.Error(w, "File path is required", http.StatusBadRequest)
        return
    }

    // Capture the current UTC timestamp
    timestamp := time.Now().UTC().String()

    // Run the Python script directly
    cmd := exec.Command("python3", "ai_model.py", filePath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Create the response
    response := Response{
        ProcessedFile: string(output),
        Timestamp:     timestamp,
    }

    // Convert response to JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the response
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func main() {
    http.HandleFunc("/process", processFile)
    fmt.Println("Starting server at port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
