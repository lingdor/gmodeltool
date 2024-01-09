package template

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/gsql"
	"github.com/lingdor/magicarray/array"
	"io"
	"io/fs"
	"reflect"
	"strings"
)

//go:embed files/schema.go.template
var schemaTmplate embed.FS

//go:embed files/new.go.template
var newtemplate embed.FS

const EndStatement = "//gmodel:gen:end"

//gmodel:gen:start:schema:tb_user
const StartStatement = "//gmodel:gen:start"

func GetNewEmptyFile(packageName string) (cnt string, err error) {
	var file fs.File
	if file, err = newtemplate.Open("files/new.go.template"); err == nil {
		defer file.Close()
		var bs []byte
		if bs, err = io.ReadAll(file); err == nil {
			cnt = string(bs)
			cnt = strings.ReplaceAll(cnt, "{$package}", packageName)
		}
	}
	return
}

func GenTableSchema(ctx context.Context, name string, db *sql.DB) (err error) {

	var version array.ZVal
	if version, err = gmodel.QueryValContext(ctx, db, gsql.Sql("select version()")); err == nil {
		var schemaArr array.MagicArray
		driverType := reflect.TypeOf(db.Driver())
		if driverType.String() == "*mysql.MySQLDriver" {
			schemaArr, err = gmodel.QueryArrRowsContext(ctx, db, gsql.Raw("desc ?", name))

		}
		//pgsql
		//schemaArr, err = gmodel.QueryArrRowsContext(ctx, db, gsql.Raw("\\d ?", name))
		if err == nil {

		}
	}
	return
	//common.LoadCommonDB()

	/*
		type TBUserEntity struct {
			Id     *int `orm:"uid,type=varchar(50)"`
			Name   *string
			IsMale bool
		}
		type TbUserSchemaType struct {
			// schema comment: userid
			Id     gmodel.Field
			Name   gmodel.Field
			IsMale gmodel.Field
		}

		var TbUserSchema = TbUserSchemaType{
			Id:   gmodel.NewField("id", gmodel.Int, 4),
			Name: gmodel.NewField("user_name", gmodel.VarChar, 50),
		}
	*/

	//fmt.Println(string(schemaTmplate))
}
