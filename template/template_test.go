package template

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/lingdor/gmodel"
	"github.com/lingdor/gmodel/gsql"
	"github.com/lingdor/gmodeltool/db"
	"github.com/lingdor/magicarray/array"
	"reflect"
	"strings"
	"testing"
)

func TestGenTableSchema(t *testing.T) {

	packageName := "xx"
	fmt.Println(GetNewEmptyFile(packageName))

}

func TestXX(t *testing.T) {
	//user:password@/dbname
	conn, err := db.Connect("mysql://root:123456@tcp(127.0.0.1:3306)/db1")
	if err != nil {
		panic(err)
	}

	driver := conn.Driver()
	tt := reflect.TypeOf(driver)
	fmt.Println(tt)
	arr, _ := gmodel.QueryArr(conn, gsql.Raw("show VARIABLES like '%version_comment%'"))
	marshal, _ := array.JsonMarshal(arr)
	fmt.Println(string(marshal))
	fmt.Println(strings.IndexAny("good", "%_"))

}

func TestXXtmp(t *testing.T) {

	cnt, err := ReadFS("files/schema.go.template")
	if err != nil {
		panic(err)
	}
	fmt.Println(cnt)
}
func TestXX2(t *testing.T) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "postgres", "123456", "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver := db.Driver()
	tt := reflect.TypeOf(driver)
	fmt.Println(tt)
	var arr array.MagicArray
	if arr, err = gmodel.QueryArr(db, gsql.Raw("select * from tb1")); err != nil {
		panic(err)
	}
	marshal, _ := array.JsonMarshal(arr)
	fmt.Println(string(marshal))

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	//user:password@/dbname
}
