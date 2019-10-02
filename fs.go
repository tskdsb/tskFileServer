package main

import (
	"net/http"
)

type SimpleFileServer struct {
	//   such as [/ . /root]
	Dir string
	// such as :8888
	Addr string
}

func (sfs *SimpleFileServer) Run() error {
	return http.ListenAndServe(sfs.Addr, http.FileServer(http.Dir(sfs.Dir)))
}
