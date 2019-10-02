package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	filePath = "filePath"
)

var (
	addr string
)

func init() {
	addr = ":9090"
}

func RunServer() {
	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		fileName := request.FormValue(filePath)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("create file error: %s\n", err)
			return
		}
		defer file.Close()
		written, err := io.Copy(file, request.Body)
		request.Body.Close()
		if err != nil {
			log.Printf("write file error: %s\n", err)
			return
		}
		log.Printf("file: %s, size: %d\n", fileName, written)
	})

	http.HandleFunc("/download", func(writer http.ResponseWriter, request *http.Request) {

	})

	
	panic(http.ListenAndServe(addr, nil))
}
