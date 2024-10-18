package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TanmayBhosale/goDB/driver"
)

func main() {
	fmt.Println("goDB server running on port: 8080 🚀")
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")

	if q == "" {
		http.Error(w, "Query not found", http.StatusPartialContent)
		return
	}

	data := driver.ProcessQuery(&q)

	resp, err2 := json.Marshal(data)

	if err2 != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}

	w.Write(resp)
}
