package gen

import (
	"github.com/lingdor/gomodeltool/gen/schema"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "gen",
	Short: "generate gmode schema code by reading database.",
	Long:  "generate gmode schema code by reading database.",
}

func init() {
	Command.AddCommand(schema.Command)

}
