package _example

import (
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/orm"
)

// step1, install command: go install github.com/lingdor/gmodeltool
// step2, edit gmodel.yml, write the right db connection dsn
// step3, chanage the below annotation parameters, --tables to your tables.
//
//go:generate gmodeltool gen schema --tables tb_user
//gmodel:gen:schema:@embed:324ae4c430a83dfac11f64246e144ff0
type TbUserSchemaType struct {
	// Id
	Id gmodel.Field
	// Name
	Name gmodel.Field
	// Age
	Age gmodel.Field
	// Createtime
	Createtime gmodel.Field
}

var TbUserSchema TbUserSchemaType = TbUserSchemaType{
	Id:         gmodel.NewField("id", "int unsigned", false, true),
	Name:       gmodel.NewField("name", "varchar(50)", true, false),
	Age:        gmodel.NewField("age", "int", true, false),
	Createtime: gmodel.NewField("createtime", "timestamp", false, false),
}

func (T *TbUserSchemaType) ToSql() (string, []any) {
	return T.TableName(), nil
}

func (T *TbUserSchemaType) As(name string) gmodel.ToSql {
	return orm.WrapField(T).As(name)
}

func (T *TbUserSchemaType) TableName() string {
	return "tb_user"
}

//gmodel:gen:end
