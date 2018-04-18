package tskdb

import (
  "encoding/json"
  "os"
  "io/ioutil"
)

type TskData struct {
  Begin int
  Size  int
  Data  interface{}
}

var (
  DataFileName = "data.db"
  AllData      map[string]TskData
)

func Save() error {
  dataFile, err := os.Create(DataFileName)
  if err != nil {
    return err
  }
  defer dataFile.Close()

  dataByte, err := json.Marshal(AllData)
  if err != nil {
    return err
  }
  _, err = dataFile.Write(dataByte)
  if err != nil {
    return err
  }

  return nil
  // var offset int
  // for name, data := range AllData {
  //   dataByte, err := json.Marshal(data)
  //   if err != nil {
  //     return err
  //   }
  //
  //   n, err := dataFile.Write(dataByte)
  //   if err != nil {
  //     return err
  //   }
  // AllData[name].Begin = offset
  // AllData[name].Size = n
  // AllData[name].Data = nil
  // offset += n
}

func Load() error {
  dataFile, err := os.Open(DataFileName)
  if err != nil {
    return err
  }
  defer dataFile.Close()

  dataByte, err := ioutil.ReadAll(dataFile)
  if err != nil {
    return err
  }

  err = json.Unmarshal(dataByte, &AllData)
  if err != nil {
    return err
  }

  return nil
}
