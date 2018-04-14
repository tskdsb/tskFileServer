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

const (
  TEMPLATE_LS string = `
<html>

<head>
  <title>File Server By TSK</title>
</head>

<body>
  {{$Path:=.Path}}
  <a href="?path={{$Path}}/..">..</a>
  <a href="?path=/">/</a>

  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Size</th>
        <th>LastModificationTime</th>
      </tr>
    </thead>

    <tbody>
      {{with .FileInfoS}}
      {{range .}}
      <tr>
        <td>
          <a href="/download?path={{$Path}}/{{.Name}}">{{.Name}}</a>
        </td>
        <td>{{humanSize .Size}}</td>
        <td>{{formatTime .ModTime}}</td>
      </tr>
      {{end}}
      {{end}}
    </tbody>
  </table>

  <br />

  <form action="/upload?path={{$Path}}" method="post" enctype="multipart/form-data">
    <input type="file" name="upload">
    <input type="submit" value="upload">
  </form>

</body>

</html>
`
)

var (
  ADDR      = ":80"
  BASE_PATH = "."

  funcMap = template.FuncMap{
    "humanSize":  humanSize,
    "formatTime": formatTime,
  }

  t = template.Must(template.New("TEMPLATE_LS").Funcs(funcMap).Parse(TEMPLATE_LS))
)

type TskHandler struct {
}

type TskFileInfo struct {
  FileInfoS []os.FileInfo
  Path      string
}

func humanSize(byteSize int64) string {
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
    fmt.Fprintf(w, "%s", BASE_PATH)
    return
  } else {
    path = filepath.Clean(path)
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
    tskFileInfo := TskFileInfo{fileInfoS, path}

    err = t.Execute(w, tskFileInfo)
    if err != nil {
      fmt.Fprintln(w, err)
    }
  } else {
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileStat.Name()))
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

  log.Printf("\nwill listen on: %s\npath to expose: %s", ADDR, BASE_PATH)
  log.Fatal(http.ListenAndServe(ADDR, nil))
}
