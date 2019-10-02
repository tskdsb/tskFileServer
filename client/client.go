package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/tskdsb/tskFileServer/util"
)

func main() {

	go func() {
		for {
			util.PrintMem()
			time.Sleep(time.Millisecond * 100)
		}
	}()

	filePath := "/Users/caicloud/git/src/github.com/tskdsb/tskssh/bin/"

	fileName := filepath.Base(filePath) + ".zip"
	url := fmt.Sprintf("http://192.168.11.100:9090/upload?filePath=%s", fileName)

	// memoryLimit := 1024 * 1024
	zipBody := bytes.NewBuffer(make([]byte, 0))

	// zipFile, err := os.Create("bin/bin.zip")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer zipFile.Close()

	fmt.Println(time.Now())
	err := util.NewZip(filePath, zipBody)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(time.Now())

	fmt.Println(zipBody.Len())
	request, err := http.NewRequest(http.MethodPut, url, zipBody)
	_, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now())
}
