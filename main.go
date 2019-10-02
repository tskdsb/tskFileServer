package main

import (
	"flag"
	"log"
)

func main() {
	var sfs SimpleFileServer

	flag.StringVar(&sfs.Dir, "dir", ".", "dir to export")
	flag.StringVar(&sfs.Addr, "addr", ".", "dir to export")
	daemon := *flag.Bool("daemon", true, "keep running")

	flag.Parse()

	if daemon {
		for {
			log.Println(sfs.Run())
		}
	} else {
		log.Println(sfs.Run())
	}
}
