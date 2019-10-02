package buffer

import (
  "io"
)

type SizedBuffer struct {
	size       int
	buf        []byte
	off        int
	readEvent  chan int
	writeEvent chan int
}

func (sb *SizedBuffer) Len() int {
	return sb.off + 1
}

func (sb *SizedBuffer) Write(p []byte) (n int, err error) {
  if sb.Len() < sb.size {

  }else {

  }
}

func (sb *SizedBuffer) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func NewSizedBuffer(size int) *SizedBuffer {
	buf := make([]byte, size, size)
	return &SizedBuffer{
		size:       size,
		buf:        buf,
		off:        0,
		readEvent:  make(chan int, 1),
		writeEvent: make(chan int, 1),
	}
}

func (sb *SizedBuffer) Size() int {
	return sb.size
}
