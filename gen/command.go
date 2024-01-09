package gen

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/gsql"
	"github.com/lingdor/gomodeltool/common"
	"github.com/lingdor/gomodeltool/config"
	"github.com/lingdor/gomodeltool/template"
	"github.com/lingdor/magicarray/array"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

var commander genSchemaCommander
var Command = &cobra.Command{
	Use:   "gen",
	Short: "generate codes for gmodel.",
	Long:  "generate codes for gmodel, example: gmodel gen schema, entity.",
	Run:   commander.runCommand,
}

type genSchemaCommander struct {
	tables      string
	name        string
	tofiles     bool
	dryRun      bool
	packageName string
}

func (g *genSchemaCommander) runCommand(cmd *cobra.Command, args []string) {

	var ctx = context.Background()
	var err error
	flags := cmd.Flags()
	flags.StringVar(&g.tables, "tables", "", "you can use comma split table names, and use wildcard character to search tables.")
	flags.StringVar(&g.name, "name", "", "")
	flags.BoolVar(&g.tofiles, "tofiles", false, "generate files, If false, generate to current file.")
	flags.BoolVar(&g.dryRun, "dry-run", false, "Testing to running and print results, do not write to files")
	common.Var(flags)

	if err = flags.Parse(args); err == nil {

		common.VerboseLog("runing generate schema progress,args: %+v", args)

		var conn *sql.DB
		if common.GetVerbose() {
			ctx = context.WithValue(ctx, gmodel.OptLogSql, true)
		}
		if err = config.LoadConfig(); err == nil {
			if conn, err = common.LoadCommonDB(); err == nil {

				sptables := strings.Split(g.tables, ",")
				for _, table := range sptables {
					if strings.ContainsAny(table, "%_") {
						var liketables array.MagicArray
						if liketables, err = gmodel.QueryArrRowsContext(ctx, conn, gsql.Raw("show tables like ?", table)); err == nil {
							iter := liketables.Iter()
							for tableRow := iter.NextVal(); tableRow != nil; tableRow = iter.NextVal() {
								tableName := tableRow.MustArr().Values().Get(0).String()
								err = g.genTableFile(ctx, tableName)
							}
						}

					} else {
						err = g.genTableFile(ctx, table)
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
	var goline = -1
	if envGoline, ok := os.LookupEnv("GOLINE"); ok {
		goline, _ = strconv.Atoi(envGoline)
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
	}
	fpathT := fmt.Sprintf("%s_t", fpath)

	var tmpFile *os.File
	if _, err = os.Stat(fpath); err == nil {
		if !g.dryRun {
			tmpFile, err = os.OpenFile(fpathT, os.O_TRUNC, 0)
		} else {
			tmpFile = os.Stdout
		}
	} else if os.IsNotExist(err) {
		tmpFile, err = g.createEmpty(ctx, fpathT)
	}
	if err == nil {
		defer func() {
			if !g.dryRun {
				tmpFile.Close()
			}
		}()
		var file *os.File
		if file, err = os.Open(fpath); err == nil {
			defer file.Close()
			posReader := common.NewPosReader(file)
			var line string
			startKey := fmt.Sprintf("%s:schema:%s", template.StartStatement, strings.ToLower(table))
			var isMatchedFirst = false
			for line, err = posReader.ReadLine(); err != nil; line, err = posReader.ReadLine() {
				if goline != 1 {
					if !isMatchedFirst && posReader.LineNo > goline {
						if strings.TrimSpace(line) == "" {
							continue //ignore empty lines
						}
						if len(line) >= len(startKey) && line[0:len(startKey)] == startKey {
							for line, err = posReader.ReadLine(); err != nil; line, err = posReader.ReadLine() {
								if len(line) >= len(template.EndStatement) && line[0:len(startKey)] == template.EndStatement {
									isMatchedFirst = true
									break
								}
							}
						}
						if err = g.genTable(ctx, table, tmpFile, startKey, template.EndStatement); err != nil {
							return err
						}
						isMatchedFirst = true
						continue
					}
					tmpFile.WriteString(line + "\n")
					continue
				}
				if !isMatchedFirst && len(line) >= len(startKey) && line[0:len(startKey)] == startKey {
					for line, err = posReader.ReadLine(); err != nil; line, err = posReader.ReadLine() {
						if len(line) >= len(template.EndStatement) && line[0:len(startKey)] == template.EndStatement {
							break
						}
					}
					if err = g.genTable(ctx, table, tmpFile, startKey, template.EndStatement); err != nil {
						return err
					}
					isMatchedFirst = true
				} else {
					_, err = tmpFile.WriteString(line + "\n")
				}
			}
		}
		if err == nil {
			if !g.dryRun {
				file.Close()
				tmpFile.Close()
				os.Remove(fpath)
				os.Rename(fpathT, fpath)
			}
		}
	}
	return
}
func (g *genSchemaCommander) createEmpty(ctx context.Context, fpath string) (w *os.File, err error) {

	pwd := common.Pwd()
	packageName := common.GetPackageName(pwd)
	var content string
	//file not exists
	common.VerboseLog("create file: %s", fpath)
	if content, err = template.GetNewEmptyFile(packageName); err == nil {
		if g.dryRun {
			w = os.Stdout
		} else {
			w, err = os.Create(fpath)
		}
		if err == nil {
			w.WriteString(content)
			w.WriteString("\n")
			return w, nil
		}
	}

	return
}

func (g *genSchemaCommander) genTable(ctx context.Context, table string, w *os.File, start, end string) error {
	common.VerboseLog("begin generate table:%s", table)
	w.WriteString(start + "\n")
	//template.GenTableSchema()
	//fmt.Printf("%+v\n", os.Environ())

	w.WriteString(end + "\n\n")
}
