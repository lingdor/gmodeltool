package template

import (
	"embed"
	_ "embed"
	"io/fs"
)

//go:embed files/schema.go.template
var schemaTmplate embed.FS

//go:embed files/new.go.template
var newtemplate embed.FS

const EndStatement = "//gmodel:gen:end"

//gmodel:gen:start:schema:tb_user
const StartStatement = "//gmodel:gen:start"

func GetNewTemplate() (fs.File, error) {
	return newtemplate.Open("files/new.go.template")
}

func GenTableSchema() {

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
