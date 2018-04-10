package tskFileServer

import (
	"net/http"
	"os"
	"fmt"
	"time"
	"html/template"
	"log"
	"io"
)

type TskHandler struct {
}

type TskFileInfo struct {
	Name    string
	Size    string
	ModTime string
}

func HumanSize(byteSize int64) string {
	var suffix = "B"
	const (
		KB int64 = 1 << ((iota + 1) * 10)
		MB
		GB
		TB
	)

	switch {
	case byteSize >= TB:
		byteSize = byteSize / TB
		suffix = "TB"
	case byteSize >= GB:
		byteSize = byteSize / GB
		suffix = "GB"
	case byteSize >= MB:
		byteSize = byteSize / MB
		suffix = "MB"
	case byteSize >= KB:
		byteSize = byteSize / KB
		suffix = "KB"
	}

	return fmt.Sprintf("%d%s", byteSize, suffix)
}

func formatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func ServeFilePath(filePath string, w http.ResponseWriter) {

	t, err := template.New("tsk").Parse(`
<head>
    <title>File Server By TSK</title>
</head>
<body>
    <a href="..">..</a>
    <a href="/">/</a>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Size</th>
                <th>LastModificationTime</th>
            </tr>
        </thead>
        {{range .}}
        <tr>
            <td>
                <a href="{{.Name}}">{{.Name}}</a>
            </td>
            <td>{{.Size}}</td>
            <td>{{.ModTime}}</td>
        </tr>
        {{end}}
    </table>
</body>
`)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	file, err := os.Open("." + filePath)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if fileStat.IsDir() {
		fileInfoS, err := file.Readdir(-1)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		TskFileInfoS := make([]TskFileInfo, len(fileInfoS))
		for i, fileInfo := range fileInfoS {
			TskFileInfoS[i].Name = fileInfo.Name()
			if fileInfo.IsDir() {
				TskFileInfoS[i].Name += "/"
			}
			TskFileInfoS[i].Size = HumanSize(fileInfo.Size())
			TskFileInfoS[i].ModTime = formatTime(fileInfo.ModTime())
		}
		t.Execute(w, TskFileInfoS)
	} else {
		_, err := io.Copy(w, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

func (h *TskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeFilePath(r.URL.Path, w)
}

func Start() {
	port := ":80"
	log.Println(port)
	log.Fatal(http.ListenAndServe(port, &TskHandler{}))
}
