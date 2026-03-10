package main

import (
	"fmt"
	"net/http"
	"os"

	"hls-player/handlers"
)

const uploadDir = "./uploads"

func main() {
	// create upload directory if it doesn't exist
	os.MkdirAll(uploadDir, os.ModePerm)

	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/upload", handlers.HandleUpload)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error:", err)
	}
}
