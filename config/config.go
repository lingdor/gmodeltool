package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var AppConfig struct {
	Gmodel struct {
		Verbose    bool
		Connection map[any]struct {
			Dsn string
		}
	}
}

func Parse(configpath string) (err error) {
	var bs []byte
	bs, err = os.ReadFile(configpath)
	if err == nil {
		err = yaml.Unmarshal(bs, &AppConfig)
	}
	return
}
func LoadConfig() (err error) {

	//GOLINE=3
	//GOPACKAGE=_example
	//GOFILE=test.go
	//PWD=/Users/zhangxiaoxu/opensource/gmodeltool/_example

	pwd, ok := os.LookupEnv("PWD")
	if !ok {
		pwd, err = os.Getwd()
	}
	if err == nil {
		return loadConfigPath(pwd, "gmodel.yml")
	}
	return
}
func loadConfigPath(pathstr string, fname string) (err error) {
	var fpath = path.Join(pathstr, fname)
	if _, err := os.Stat(fpath); err == nil {
		return Parse(fpath)
	}
	pathstr = path.Dir(pathstr)
	if pathstr == "/" {
		return fmt.Errorf("no found config file: gmodel.yml")
	}
	return loadConfigPath(pathstr, fname)
}
