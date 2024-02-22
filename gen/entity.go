package gen

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lingdor/gmodeltool/common"
	"github.com/lingdor/gmodeltool/template"
	"github.com/lingdor/magicarray/utils"
	"strings"
)

func (g *genSchemaCommander) GenTableEntity(ctx context.Context, tname string, name string, columns []*common.ColumnInfo) (code string, imports []common.Import, err error) {

	if len(columns) < 1 {
		return "", []common.Import{}, nil
	}

	imports = make([]common.Import, 0)
	imports = append(imports, common.Import{Path: "github.com/lingdor/gmodel"})
	maxLen := g.maxColumnLen(columns)
	maxLen += 8
	structBuf := &bytes.Buffer{}
	casesBuf := &bytes.Buffer{}
	var typeName = fmt.Sprintf("%sEntity", name)
	structBuf.WriteString(fmt.Sprintf("type %s struct{\n", typeName))

	for _, column := range columns {
		field := column.Field

		fillName := utils.PadLeftRightSpaces(column.Name, 4, maxLen)

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
			if g.parseTime {
				memberType = "*time.Time"
				imports = append(imports, common.Import{Path: "time"})
			}
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
		structBuf.WriteString(fmt.Sprintf("    // %s %s\n", column.Name, strings.ReplaceAll(column.Comment, "\n", "\n    //")))
		structBuf.WriteString(fmt.Sprintf("%s %s `%s`\n", fillName, memberType, tagInfo))

		casesBuf.WriteString(fmt.Sprintf("    case %q: handlers[i]=&entity.%s\n", field.Name(), column.Name))

	}
	structBuf.WriteString("\n}")

	if code, err = template.ReadFS("files/entity.go.template"); err == nil {
		code = strings.ReplaceAll(code, "${struct}", structBuf.String())
		code = strings.ReplaceAll(code, "${structName}", typeName)
		code = strings.ReplaceAll(code, "${cases}", casesBuf.String())
		code = strings.ReplaceAll(code, "${tableName}", fmt.Sprintf("%q", tname))
	}
	return
}
