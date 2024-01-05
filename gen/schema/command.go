package schema

import (
	"fmt"
	"github.com/lingdor/gomodeltool/schema"
	"github.com/spf13/cobra"
	"os"
)

var Command = &cobra.Command{
	Use:   "schema",
	Short: "generate gmode schema code by reading database.",
	Long:  "generate gmode schema code by reading database.",
	Run:   runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {

	/*
		1, load config
		2, connect db
		3, read schema
		4, generate files / current file -> memory
		5, import write
		6, template replace
		7,

	*/

	schema.GenTableSchema()
	fmt.Printf("%+v\n", os.Environ())
	//
	//path:=os.LookupEnv("pwd")
	//cmd.Flags().

}
