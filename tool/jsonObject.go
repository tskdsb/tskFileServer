package tool

import (
  "encoding/json"
)

type JsonObject map[string]interface{}

func (j JsonObject) ToString() (s string, err error) {
  data, err := json.Marshal(j)
  if err != nil {
    return
  }

  s = string(data)
  return
}
