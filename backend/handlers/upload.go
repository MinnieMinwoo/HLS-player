package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"
const maxUploadSize = 500 * 1024 * 1024 // 500MB

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// restrict file size
	r.ParseMultipartForm(maxUploadSize)

	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "File upload failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// check file extension
	ext := filepath.Ext(handler.Filename)
	validExtensions := map[string]bool{
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
		".flv": true, ".wmv": true, ".webm": true, ".m3u8": true,
	}

	if !validExtensions[ext] {
		http.Error(w, "Invalid file format. Allowed: mp4, mov, avi, mkv, flv, wmv, webm", http.StatusBadRequest)
		return
	}

	// delete existing files (keep only the latest)
	dir, err := os.Open(uploadDir)
	if err == nil {
		defer dir.Close()
		files, _ := dir.Readdirnames(-1)
		for _, f := range files {
			os.Remove(filepath.Join(uploadDir, f))
		}
	}

	// save new file
	filePath := filepath.Join(uploadDir, "latest"+ext)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to write file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success","message":"Video uploaded successfully","filename":"latest%s"}`, ext)
}
