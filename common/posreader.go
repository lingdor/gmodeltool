package common

import (
	"bufio"
	"io"
	"os"
)

type PosReader struct {
	Reader *bufio.Reader
	File   *os.File
	//StartPos int64
	LineNo int
}

func NewPosReader(file *os.File) *PosReader {
	if file == nil {
		return &PosReader{Reader: nil}
	}
	return &PosReader{
		Reader: bufio.NewReader(file),
	}
}
func (r *PosReader) ReadLine() (line string, err error) {
	if r.Reader == nil {
		return "", io.EOF
	}
	//if r.StartPos, err = r.File.Seek(0, 1); err == nil {
	var bs []byte
	if bs, _, err = r.Reader.ReadLine(); err == nil {
		r.LineNo++
		line = string(bs)
	}
	//}
	return
}
