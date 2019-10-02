package util

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const Ki = 1 << 10
const Mi = 1 << 20

func PrintMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d\n", m.Sys/Mi)
}

func NewZip(filePath string, writer io.Writer) error {
	newWriter := zip.NewWriter(writer)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		filePath = EnsureDirPath(filePath)

		err := filepath.Walk(filepath.Dir(filePath), func(path string, info os.FileInfo, err error) error {

			if info.IsDir() {
				dirName := path

				dirName, err = filepath.Rel(filePath, dirName)
				if err != nil {
					return err
				}
				if dirName == "." {
					return nil
				}

				// 如果末尾不带/，文件夹会被当成文件压入包内，解压时也会是一个文件
				dirName = EnsureDirPath(dirName)
				fmt.Println(dirName)
				_, err = newWriter.Create(dirName)
				if err != nil {
					log.Printf("create dir(%s) error: %s", dirName, err)
					return err
				}
			} else {
				dirName, err := filepath.Rel(filePath, path)
				if err != nil {
					return err
				}
				fmt.Println(dirName)
				create, err := newWriter.Create(dirName)
				if err != nil {
					log.Printf("create file(%s) error: %s", info.Name(), err)
					return err
				}
				file, err := os.Open(path)
				if err != nil {
					log.Printf("open file(%s) error: %s", path, err)
					return err
				}
				defer file.Close()

				wn, err := io.Copy(create, file)
				if err != nil {
					log.Printf("copy file(%s) error: %s", info.Name(), err)
					return err
				}
				log.Printf("file: %s, size: %d", file.Name(), wn)
			}

			return nil
		})

		if err != nil {
			log.Printf("err after walk: %s", err)
			return err
		}
	} else {
		create, err := newWriter.Create(info.Name())
		if err != nil {
			return err
		}
		_, err = io.Copy(create, file)
		if err != nil {
			return err
		}
	}

	err = newWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// 如果是文件夹，末尾加/
func EnsureDirPath(dirPath string) string {
	if dirPath[len(dirPath)-1] != '/' {
		dirPath += "/"
	}

	return dirPath
}

func SendMultipart(url string, field string, name string) (*http.Response, error) {
	r, w := io.Pipe()
	mWriter := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer mWriter.Close()

		fWriter, err := mWriter.CreateFormFile(field, name)
		if err != nil {
			return
		}

		file, err := os.Open(name)
		if err != nil {
			return
		}
		defer file.Close()

		if _, err = io.Copy(fWriter, file); err != nil {
			return
		}
	}()

	return http.Post(url, mWriter.FormDataContentType(), r)
}

func ReceiveMultipart(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// demonstration purpose - we have 2 fields named: text_field and file_file
	// they will be parsed in order
	var _ = []string{"text_field", "file_field"}

	// parse text field
	var text = make([]byte, 512)
	p, err := reader.NextPart()
	// one more field to parse, EOF is considered as failure here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "text_field" {
		http.Error(w, "text_field is expected", http.StatusBadRequest)
	}

	_, err = p.Read(text)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse file field
	p, err = reader.NextPart()
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "file_field" {
		http.Error(w, "file_field is expected", http.StatusBadRequest)
	}

	buf := bufio.NewReader(p)
	sniff, _ := buf.Peek(512)
	contentType := http.DetectContentType(sniff)
	if contentType != "application/zip" {
		http.Error(w, "file type not allowed", http.StatusBadRequest)
		return
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	var maxSize int64 = 32 << 20
	lmt := io.LimitReader(p, maxSize+1)
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if written > maxSize {
		os.Remove(f.Name())
		http.Error(w, "file size over limit", http.StatusBadRequest)
		return
	}
}
