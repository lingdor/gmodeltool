package template

import (
	"embed"
	_ "embed"
	"io"
	"io/fs"
	"strings"
)

//go:embed files/*
var templateFS embed.FS

const EndStatement = "//gmodel:gen:end"

//gmodel:gen:start:schema:tb_user
const StartStatement = "//gmodel:gen"

func GetNewEmptyFile(packageName string) (cnt string, err error) {
	cnt, err = ReadFS("files/new.go.template")
	cnt = strings.ReplaceAll(cnt, "{$package}", packageName)
	return
}
func ReadFS(fname string) (cnt string, err error) {
	var file fs.File
	if file, err = templateFS.Open(fname); err == nil {
		defer file.Close()
		var bs []byte
		if bs, err = io.ReadAll(file); err == nil {
			cnt = string(bs)
		}
	}
	return
}
