package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/lingdor/gmodeltool/common"
	"github.com/lingdor/gmodeltool/config"
	"github.com/lingdor/gmodeltool/gen"
	"github.com/spf13/cobra"
)

func main() {

	var rootCommand = &cobra.Command{
		Use:   "gmodeltool",
		Short: "gmodeltool used for generate gmodel codes",
		Long:  "gmodeltool used for generate gmodel codes.",
	}

	var err error
	if err = config.LoadConfig(); err == nil {
		common.InitCommand(rootCommand)
		rootCommand.AddCommand(gen.Command)
		err = rootCommand.Execute()
	}
	if err != nil {
		panic(err)
	}
}
