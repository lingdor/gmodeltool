package main

import (
	"github.com/lingdor/gomodeltool/config"
	"github.com/lingdor/gomodeltool/gen"
)

func main() {
	var err error
	if err = config.LoadConfig(); err == nil {
		err = gen.Command.Execute()
	}
	if err != nil {
		panic(err)
	}
}
