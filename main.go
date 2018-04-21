package main

import (
  "net/http"
  "log"
  "tskFileServer/fileServer"
  "flag"
)

func main() {

  flag.StringVar(&fileServer.ADDR, "addr", "", "[ip:port], such as: 0.0.0.0:80")
  flag.StringVar(&fileServer.BASE_PATH, "dir", ".", "Directory to expose, such as: /tmp/share")
  flag.Parse()

  http.HandleFunc("/", fileServer.RouterHandler)
  http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(fileServer.BASE_PATH))))

  http.HandleFunc("/download", fileServer.Download)
  http.HandleFunc("/upload", fileServer.Upload)
  http.HandleFunc("/mkdir", fileServer.Mkdir)

  http.HandleFunc("/cmd/local", fileServer.CmdLocal)
  http.HandleFunc("/cmd/ssh", fileServer.CmdSsh)

  log.Printf("\nwill listen on: %s\npath to expose: %s", fileServer.ADDR, fileServer.BASE_PATH)
  log.Fatal(http.ListenAndServe(fileServer.ADDR, nil))
}
