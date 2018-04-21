package fileServer

import (
  "net/http"
  "os"
  "fmt"
  "time"
  "html/template"
  "io"
  "path/filepath"
  "reflect"
  "io/ioutil"
  "tskFileServer/tool"
  "tskFileServer/cmd"
)

const (
  TEMPLATE_LS string = `
<html>

<head>
  <title>File Server By TSK</title>
</head>

<body>
  {{$Path:=.Path}}

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
      <tr>
        <td>
          <a href="?path={{$Path}}/..">..</a>
        </td>
      </tr>
      {{with .FileInfoS}}
      {{range .}}
      <tr>
        <td>
          <a href="/download?path={{$Path}}/{{.Name}}&disposition=inline">{{dirSuffix .}}</a>
        </td>
        <td>{{humanSize .Size}}</td>
        <td>{{formatTime .ModTime}}</td>
        {{if not (isDir .)}}
        <td>
          <a href="/download?path={{$Path}}/{{.Name}}&disposition=attachment">download</a>
        </td>
        {{end}}
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
  
  <form action="/mkdir?path={{$Path}}" method="post" enctype="multipart/form-data">
    <input type="text" name="mkdir">
    <input type="submit" value="mkdir">
  </form>

</body>

</html>
`
)

var (
  ADDR      string
  BASE_PATH string

  funcMap = template.FuncMap{
    "humanSize":  humanSize,
    "formatTime": formatTime,
    "dirSuffix":  dirSuffix,
    "isDir":      isDir,
  }

  t = template.Must(template.New("TEMPLATE_LS").Funcs(funcMap).Parse(TEMPLATE_LS))
)

type TskHandler struct {
}

type TskFileInfo struct {
  FileInfoS []os.FileInfo
  Path      string
}

func dirSuffix(f os.FileInfo) string {
  if f.IsDir() {
    return f.Name() + "/"
  } else {
    return f.Name()
  }
}
func isDir(f os.FileInfo) bool {
  return f.IsDir()
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

func Download(w http.ResponseWriter, r *http.Request) {
  path := filepath.Clean(r.FormValue("path"))

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
    // w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // w.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
    disposition := r.FormValue("disposition")
    w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, fileStat.Name()))
    w.Header().Set("Content-Length", fmt.Sprintf("%s", fileStat.Size()))
    _, err := io.Copy(w, file)
    if err != nil {
      fmt.Fprintln(w, err)
    }
  }
}

func Upload(w http.ResponseWriter, r *http.Request) {
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
    return
  }

  http.Redirect(w, r, "/download?path="+path, http.StatusTemporaryRedirect)
}

func Mkdir(w http.ResponseWriter, r *http.Request) {
  path := r.FormValue("path")
  dir := r.FormValue("mkdir")
  dirToMk := filepath.Join(BASE_PATH, path, dir)
  err := os.MkdirAll(dirToMk, os.ModePerm)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  http.Redirect(w, r, "/download?path="+path, http.StatusTemporaryRedirect)
}

func RouterHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    mKeyS := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m").MapKeys()
    for p := range mKeyS {
      fmt.Fprintf(w, "<a href=\"%s\">%s</a><br />", mKeyS[p], mKeyS[p])
    }
  }
}

func CmdLocal(w http.ResponseWriter, r *http.Request) {
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }
  defer r.Body.Close()

  var cmdObject = cmd.CmdObject{}
  err = tool.Json2object(data, &cmdObject)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  resultJson, err := tool.Object2json(cmd.RunLocal(cmdObject))
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  fmt.Fprintf(w, "%s", resultJson)
}

func CmdSsh(w http.ResponseWriter, r *http.Request) {
  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }
  defer r.Body.Close()

  var cmdObject = cmd.SshCmdObject{}
  err = tool.Json2object(data, &cmdObject)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  resultCmd := cmd.RunSsh(cmdObject)

  resultJson, err := tool.Object2json(resultCmd)
  if err != nil {
    fmt.Fprintln(w, err)
    return
  }

  fmt.Fprintf(w, "%s", resultJson)
}

func (h *TskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
