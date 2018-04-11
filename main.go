package tskFileServer

import (
	"net/http"
	"os"
	"fmt"
	"time"
	"html/template"
	"log"
	"io"
	"path/filepath"
)

var (
	PORT      = ":80"
	BASE_PATH = "./"
)

type TskHandler struct {
}

type TskFileInfo struct {
	Name    string
	Size    string
	ModTime string
	Path    string
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

func download(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if path == "" {
		path = "."
	}
	t, err := template.New("tsk").Parse(`
<head>
  <title>File Server By TSK</title>
</head>

<body>
  <a href="?path=..">..</a>
  <a href="?path=/">/</a>
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
        <a href="/download?path={{.Path}}{{.Name}}">{{.Name}}</a>
      </td>
      <td>{{.Size}}</td>
      <td>{{.ModTime}}</td>
    </tr>
    {{end}}
  </table>
  <form action="/upload?path={{.Path}}" method="POST" enctype="multipart/form-data">
    <input type="file" name="upload">
    <input type="submit" value="upload">
  </form>
</body>
`)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	file, err := os.Open(filepath.Join(BASE_PATH, path))
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
			TskFileInfoS[i].Path = path
			if path[len(path)-1] != '/' {
				TskFileInfoS[i].Path += "/"
			}
		}
		t.Execute(w, TskFileInfoS)
	} else {
		_, err := io.Copy(w, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	file, fileHeader, err := r.FormFile("upload")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	file2, err := os.Create(filepath.Join(BASE_PATH, path, fileHeader.Filename))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file2.Close()

	_, err = io.Copy(file2, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	http.Redirect(w, r, "/download?path="+path, http.StatusTemporaryRedirect)
}

func (h *TskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func Start() {

	//show routers
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	var mKeyS = reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m").MapKeys()
	//	for p := range mKeyS {
	//		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br />", mKeyS[p], mKeyS[p])
	//	}
	//})

	http.HandleFunc("/download", download)
	http.HandleFunc("/upload", upload)

	log.Println(PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
