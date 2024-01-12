package gen

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lingdor/gmodeltool/common"
	"github.com/lingdor/gmodeltool/template"
	"os"
	"strings"
)

func (g *genSchemaCommander) GenEntity(ctx context.Context, w *os.File, table, objName string, fields []*common.ColumnInfo, start, end string) {
	if code, err := g.GenTableEntity(ctx, table, objName, fields); err == nil {

		w.WriteString(fmt.Sprintf("%s:%s\n", start, common.MD5(code)))
		w.WriteString(code)
		w.WriteString(fmt.Sprintf("\n%s", end))
	}

}

func (g *genSchemaCommander) GenTableEntity(ctx context.Context, tname string, name string, columns []*common.ColumnInfo) (code string, err error) {

	structBuf := bytes.Buffer{}
	var typeName = fmt.Sprintf("%sEntity", name)
	structBuf.WriteString(fmt.Sprintf("type %s struct{\n", typeName))
	for _, column := range columns {
		field := column.Field
		memberType := "*string"
		typeStr := strings.ToLower(field.Type())
		if index := strings.Index(typeStr, "("); index > -1 {
			typeStr = typeStr[0:strings.Index(typeStr, "(")]
		}
		switch typeStr {
		case "int":
			memberType = "*int"
		case "tinyint":
			memberType = "*int8"
		case "mediumint":
			memberType = "*int16"
		case "bigint":
			memberType = "*int64"
		case "tinyint unsigned":
			memberType = "*uint8"
		case "mediumint unsigned":
			memberType = "*uint16"
		case "int tinyint":
			memberType = "*uint"
		case "bigint bigint":
			memberType = "*int64"
		case "date", "timestamp", "time", "datetime":
			memberType = "*time.Time"
		case "float":
			memberType = "*float"
		case "double":
			memberType = "*double"
		}
		tagInfo := fmt.Sprintf(`gmodel:"%s"`, field.Name())
		if g.gorm {
			ops := ""
			if field.IsPK() {
				ops = ";primaryKey;"
			}
			tagInfo = fmt.Sprintf(`%s gorm:"column:%s%s"`, tagInfo, field.Name(), ops)
		}
		structBuf.WriteString(fmt.Sprintf("    %s	%s `%s` //%s\n", column.Name, memberType, tagInfo, column.Comment))
	}
	structBuf.WriteString("\n}")

	if code, err = template.ReadFS("files/entity.go.template"); err == nil {
		code = strings.ReplaceAll(code, "{$struct}", structBuf.String())
	}
	return
}
