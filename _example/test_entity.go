package _example

import "time"

// TbUserEntity
// step1, install command: go install github.com/lingdor/gmodeltool
// step2, edit gmodel.yml, write the right db connection dsn
// step3, chanage the below annotation parameters, --tables to your tables.
//
//go:generate gmodeltool gen entity --gorm --tables tb_user
//gmodel:gen:entity:@embed:c93b35b26ef9cdaf339981ccb89bcd3c
type TbUserEntity struct {
	// Id
	Id *string `gmodel:"id" gorm:"column:id;primaryKey;"` //
	// Name
	Name *string `gmodel:"name" gorm:"column:name"` //
	// Age
	Age *int `gmodel:"age" gorm:"column:age"` //
	// Createtime
	Createtime *time.Time `gmodel:"createtime" gorm:"column:createtime"` //

}

//func (entity *TbUserEntity) GetFieldsHandler(fields []string)[]any{
//    handlers :=make([]any,len(fields))
//    lack := false;
//    for i,field:=range fields {
//        switch(field){
//            case "id":handlers[i]=&t.Id
//            case "name":handlers[i]=&t.Name
//            default:lack=true
//        }
//    }
//    if lack{
//        handlers=gmodel.fillHandlers(t, fields, handlers)
//    }
//    return handlers
//}
//

//}
//gmodel:gen:end
