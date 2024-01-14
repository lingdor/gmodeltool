package gen

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/gsql"
	"github.com/lingdor/gmodeltool/common"
	"github.com/lingdor/gmodeltool/config"
	"github.com/lingdor/gmodeltool/template"
	"github.com/lingdor/magicarray/array"
	"github.com/spf13/cobra"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var commander genSchemaCommander
var Command = &cobra.Command{
	Use:   "gen",
	Short: "generate codes for gmodel, example: gmodel gen schema, entity",
	//ValidArgs: []string{"tables"},
	Example: "gmodeltool gen --tables='tb_%'",
	Long:    "generate codes for gmodel, example: gmodel gen schema, entity.",
	Run:     commander.runCommand,
}

func init() {
	flags := Command.Flags()
	flags.StringVar(&commander.tables, "tables", "", "you can use comma split table names, and use wildcard character to search tables.")
	flags.StringVar(&commander.name, "name", "", "")
	flags.StringVarP(&commander.tofiles, "to-files", "t", "", "generate files to path, If empty, generate to current file.")
	flags.BoolVar(&commander.dryRun, "dry-run", false, "Testing to running and print results, do not write to files")
	flags.BoolVar(&commander.gorm, "gorm", false, "generate entity for gorm tags.")
}

type genSchemaCommander struct {
	tables      string
	name        string
	tofiles     string
	dryRun      bool
	packageName string
	action      string
	gorm        bool
}

const CommandSchema = "schema"
const CommandEntity = "entity"

func (g *genSchemaCommander) runCommand(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("no found gen commands like: schema,entity, for example: gmodeltool gen schema --tables=tb1")
		return
	}
	g.action = strings.ToLower(strings.TrimSpace(args[0]))
	if g.action != CommandSchema && g.action != CommandEntity {
		fmt.Println("no found gen commands like: schema,entity, for example: gmodeltool gen schema --tables=tb1")
		return
	}

	var ctx = context.Background()
	var err error
	//if err = flags.Parse(args); err == nil {

	common.VerboseLog("runing generate schema progress,args: %+v", args)
	common.VerboseLog("param tables:%s", g.tables)

	var conn *sql.DB
	if common.GetVerbose() {
		ctx = context.WithValue(ctx, gmodel.OptLogSql, true)
	}
	if err = config.LoadConfig(); err == nil {
		common.VerboseLog("config loaded success")
		if conn, err = common.LoadCommonDB(); err == nil {

			tables := make([]string, 0, 10)
			sptables := strings.Split(g.tables, ",")
			for _, table := range sptables {
				if strings.TrimSpace(table) == "" {
					continue
				}
				if strings.ContainsAny(table, "%_") {
					var liketables array.MagicArray
					if liketables, err = gmodel.QueryArrRowsContext(ctx, conn, gsql.Sql(fmt.Sprintf("show tables like '%s'", strings.ReplaceAll(table, "'", "''")))); err == nil {
						if array.Empty(liketables) {
							panic(fmt.Errorf("no found tables:%s", g.tables))
						}
						iter := liketables.Iter()
						for tableRow := iter.FirstVal(); tableRow != nil; tableRow = iter.NextVal() {
							tableName := tableRow.MustArr().Values().Get(0).String()
							tables = append(tables, tableName)
							//err = g.genTableFile(ctx, conn, tableName)
						}
					}

				} else {
					tables = append(tables, table)
					//err = g.genTableFile(ctx, conn, table)
				}
				g.genTables(ctx, conn, tables)
			}
		}
	}
	//}
	if err != nil {
		panic(err)
	}
}

func (g *genSchemaCommander) GenName(dbName string) string {
	if strings.Index(dbName, "_") != -1 {
		arr := array.NewArr(dbName)
		arr = array.WashAll(arr, array.GetWashFuncWashUnderScoreCaseToCamelCase(true))
		dbName = arr.Get(0).String()
	} else if len(dbName) > 0 {
		runes := []rune(dbName)
		runes[0] = unicode.ToUpper(runes[0])
		dbName = string(runes)
	} else {
		panic("db name is empty")
	}
	return dbName
}

func (g *genSchemaCommander) genTables(ctx context.Context, conn *sql.DB, tables []string) (err error) {

	common.VerboseLog("generating table: %v", tables)
	pwd := common.Pwd()
	var fpath string

	if gofile, ok := os.LookupEnv("GOFILE"); !ok && g.tofiles == "" {
		return fmt.Errorf("no found --to-files parameters set")
	} else if ok && g.tofiles == "" {
		fpath = path.Join(pwd, gofile)
	}

	if g.tofiles != "" {
		for _, table := range tables {
			fname := strings.ToLower(strings.ReplaceAll(table, "_", ""))
			fname = fmt.Sprintf("%s_gen.go", fname)
			var to string
			if g.tofiles == "" || g.tofiles[0] == '.' {
				to = path.Join(pwd, g.tofiles)
			} else {
				to = g.tofiles
			}
			fpath = path.Join(to, fname)
			startKey := fmt.Sprintf("%s:%s:%s", template.StartStatement, g.action, strings.ToLower(table))
			if err = g.genTableFile(ctx, conn, fpath, startKey, table); err != nil {
				return err
			}
		}
		return
	}
	startKey := fmt.Sprintf("%s:%s:@embed", template.StartStatement, g.action)
	err = g.genTableFile(ctx, conn, fpath, startKey, tables...)
	return
}
func (g *genSchemaCommander) genTableFile(ctx context.Context, conn *sql.DB, fpath string, startKey string, tables ...string) (err error) {

	var goline = -1
	if envGoline, ok := os.LookupEnv("GOLINE"); ok {
		goline, _ = strconv.Atoi(envGoline)
	}
	fpathT := fmt.Sprintf("%s_t", fpath)

	var tmpFile *os.File
	var file *os.File
	if _, err = os.Stat(fpath); err == nil {
		file, err = os.Open(fpath)
		if !g.dryRun {
			tmpFile, err = os.OpenFile(fpathT, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
		} else {
			tmpFile = os.Stdout
		}
	} else if os.IsNotExist(err) {
		tmpFile, err = g.createEmpty(ctx, fpathT)
		file = nil
	}
	if err == nil {
		defer func() {
			if !g.dryRun {
				tmpFile.Close()
			}
		}()
		if err == nil {
			defer file.Close()
			posReader := common.NewPosReader(file)

			var isMatchedFirst = false
			if file != nil {
				for line, err := posReader.ReadLine(); err == nil; line, err = posReader.ReadLine() {
					if goline != -1 && g.tofiles == "" { //embed model
						if !isMatchedFirst && posReader.LineNo > goline {
							if strings.TrimSpace(line) == "" {
								continue //ignore empty lines
							}
							if len(line) >= len(startKey) && strings.ToLower(line[0:len(startKey)]) == startKey {
								for line, err = posReader.ReadLine(); err == nil; line, err = posReader.ReadLine() {
									if len(line) >= len(template.EndStatement) && line[0:len(template.EndStatement)] == template.EndStatement {
										isMatchedFirst = true
										break
									}
								}
							}
							if err = g.genTable(ctx, conn, tables, tmpFile, startKey, template.EndStatement); err != nil {
								return err
							}
							isMatchedFirst = true
							continue
						}
						tmpFile.WriteString(line + "\n")
						continue
					}
					if !isMatchedFirst && len(line) >= len(startKey) && line[0:len(startKey)] == startKey {
						for line, err = posReader.ReadLine(); err == nil; line, err = posReader.ReadLine() {
							if len(line) >= len(template.EndStatement) && strings.ToLower(line[0:len(template.EndStatement)]) == template.EndStatement {
								break
							}
						}
						if err = g.genTable(ctx, conn, tables, tmpFile, startKey, template.EndStatement); err != nil {
							return err
						}
						isMatchedFirst = true
					} else {
						_, err = tmpFile.WriteString(line + "\n")
					}
				}
			}
			if !isMatchedFirst {
				err = g.genTable(ctx, conn, tables, tmpFile, startKey, template.EndStatement)
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

func (g *genSchemaCommander) genTable(ctx context.Context, conn *sql.DB, tables []string, w *os.File, start, end string) (err error) {
	common.VerboseLog("begin generate table:%v", tables)

	buff := bytes.Buffer{}
	for _, table := range tables {
		//var fields = make([]*gmodel.FieldInfo, 0, 10)
		//fields := array.Make(true, true, 10)
		fields := make([]*common.ColumnInfo, 0, 10)

		driverType := reflect.TypeOf(conn.Driver())
		if strings.Contains(driverType.String(), "mysql") {
			var descArr array.MagicArray

			if descArr, err = gmodel.QueryArrRowsContext(ctx, conn, gsql.Sql(fmt.Sprintf("show full fields from `%s`", table))); err == nil {
				iter := descArr.Iter()
				for row := iter.FirstVal(); row != nil; row = iter.NextVal() {
					if row, ok := row.Arr(); ok {
						cField := row.Get("Field").String()
						cType := row.Get("Type").String()
						cNull := row.Get("Null").String()
						cKey := row.Get("Key").String()
						cComment := row.Get("Comment").String()
						field := gmodel.NewField(cField, cType, cNull == "YES", cKey == "PRI")
						fields = append(fields, &common.ColumnInfo{
							Field:   field,
							Comment: cComment,
							Name:    g.GenName(cField),
						})
					}
				}
			}

		} else if strings.Contains(driverType.String(), "pq.Driver") {
			//pgsql todo
			return fmt.Errorf("not supported pgsql yet")
		}
		ObjName := g.GenName(table)
		var code string
		if g.action == CommandSchema {
			code, err = g.GenTableSchema(ctx, table, ObjName, fields)
		} else if g.action == CommandEntity {
			code, err = g.GenTableEntity(ctx, table, ObjName, fields)
		} else {
			return fmt.Errorf("no found command:%s", g.action)
		}
		if err == nil {
			buff.WriteString(code)
			buff.WriteString("\n")
		}
	}
	data := buff.Bytes()
	w.WriteString(fmt.Sprintf("%s:%s\n", start, common.MD5(data)))
	w.Write(data)
	w.WriteString(fmt.Sprintf("%s\n", end))
	return
}
func (g *genSchemaCommander) maxColumnLen(cols []*common.ColumnInfo) int {

	max := 0
	for _, col := range cols {
		if len(col.Name) > max {
			max = len(col.Name)
		}
	}
	return max

}
