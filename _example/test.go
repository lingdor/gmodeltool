package _example

import "time"

// step1, install command: go install github.com/lingdor/gmodeltool
// step2, edit gmodel.yml, write the right db connection dsn
// step3, chanage the below annotation parameters, --tables to your tables.
// go :generate gmodeltool gen schema --tables t_user
// or
//
//go:generate gmodeltool gen entity --gorm --tables tb_user
//gmodel:gen:entity:@embed:10ee7e5ec910d29e251a7a7481b7fed9
type TbUserEntity struct {
	id         *string    `gmodel:"id" gorm:"column:id;primaryKey;"`     //
	name       *string    `gmodel:"name" gorm:"column:name"`             //
	age        *int       `gmodel:"age" gorm:"column:age"`               //
	createtime *time.Time `gmodel:"createtime" gorm:"column:createtime"` //

}

//gmodel:gen:end
