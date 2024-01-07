package schema

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/gsql"
	"github.com/lingdor/gomodeltool/common"
	"github.com/lingdor/gomodeltool/config"
	"github.com/lingdor/gomodeltool/db"
	"github.com/lingdor/gomodeltool/log"
	"github.com/lingdor/gomodeltool/template"
	"github.com/lingdor/magicarray/array"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"unicode"
)

var commander genSchemaCommander
var Command = &cobra.Command{
	Use:   "schema",
	Short: "generate gmode schema code by reading database.",
	Long:  "generate gmode schema code by reading database.",
	Run:   commander.runCommand,
}

type genSchemaCommander struct {
	tables  string
	name    string
	tofiles bool
	dryRun  bool
}

func (g *genSchemaCommander) runCommand(cmd *cobra.Command, args []string) {

	var ctx = context.Background()
	var err error
	flags := cmd.Flags()
	flags.StringVar(&g.tables, "tables", "", "you can use comma split table names, and use wildcard character to search tables.")
	flags.StringVar(&g.name, "name", "", "")
	flags.BoolVar(&g.tofiles, "tofiles", false, "generate files, If false, generate to current file.")
	flags.BoolVar(&g.dryRun, "dry-run", false, "Testing to running and print results, do not write to files")
	db.Var(flags)
	log.Var(flags)
	if err = flags.Parse(args); err == nil {

		log.VerboseLog("runing generate schema progress,args: %+v", args)

		var conn *sql.DB
		if log.GetVerbose() {
			ctx = context.WithValue(ctx, gmodel.OptLogSql, true)
		}
		if err = config.LoadConfig(); err == nil {
			if conn, err = db.LoadCommonDB(); err == nil {

				sptables := strings.Split(g.tables, ",")
				for _, table := range sptables {
					if strings.ContainsAny(table, "%_") {
						var liketables array.MagicArray
						if liketables, err = gmodel.QueryArrRowsContext(ctx, conn, gsql.Raw("show tables like ?", table)); err == nil {
							iter := liketables.Iter()
							for tableRow := iter.NextVal(); tableRow != nil; tableRow = iter.NextVal() {
								tableName := tableRow.MustArr().Values().Get(0).String()
								err = g.genTable(ctx, tableName)
							}
						}

					} else {
						err = g.genTable(ctx, table)
					}
				}
			}
		}
	}
	if err != nil {
		panic(err)
	}
}

func (g *genSchemaCommander) genTableFile(ctx context.Context, table string) (err error) {

	pwd := common.Pwd()
	name := g.name
	var fpath string
	if strings.Index(name, "_") != -1 {
		arr := array.NewArr(name)
		arr = array.WashAll(arr, array.GetWashFuncWashUnderScoreCaseToCamelCase(true))
		name = arr.Get(0).String()
	} else if len(name) > 0 {
		runes := []rune(name)
		runes[0] = unicode.ToUpper(runes[0])
		name = string(runes)
	} else {
		return fmt.Errorf("table name is empty")
	}

	if gofile, ok := os.LookupEnv("GOFILE"); !ok {
		g.tofiles = true
	} else {
		fpath = path.Join(pwd, gofile)
	}

	if g.tofiles {
		fname := strings.ToLower(strings.ReplaceAll(table, "_", ""))
		fname = fmt.Sprintf("%s_gen.go", fname)
		fpath = path.Join(pwd, fname)
		if err = g.checkCreate(ctx, fpath, pwd); err != nil {
			return err
		}
	}
	if g.dryRun {
		return g.genTable(ctx, table)
	}

	//check start
	// check end

	//write

}
func (g *genSchemaCommander) checkCreate(ctx context.Context, fpath string, pwd string) (err error) {

	var packageName string
	if envPackage, ok := os.LookupEnv("GOPACKAGE"); !ok {
		if pwd[len(pwd)-1] == '/' {
			pwd = pwd[0 : len(pwd)-1]
		}
		index := strings.LastIndex(pwd, string(os.PathSeparator))
		if index != -1 {
			packageName = pwd[index+1:]
		} else {
			packageName = "main"
		}
	} else {
		packageName = envPackage
	}
	if _, err = os.Stat(fpath); err != nil && os.IsNotExist(err) {
		//file not exists
		var file *os.File
		var templateFile fs.File
		common.VerboseLog("create file: %s", fpath)
		if templateFile, err = template.GetNewTemplate(); err == nil {
			defer templateFile.Close()
			var newBS []byte
			if newBS, err = io.ReadAll(templateFile); err == nil {
				content := string(newBS)
				content = strings.ReplaceAll(content, "{$package}", packageName)
				if g.dryRun {
					fmt.Println(content)
					return
				}
				if file, err = os.Create(fpath); err == nil {
					defer file.Close()
					_, err = file.WriteString(content)
				}
			}
		}
	}
	return
}

func (g *genSchemaCommander) genTable(ctx context.Context, table string) error {

	template.GenTableSchema()
	fmt.Printf("%+v\n", os.Environ())
}
