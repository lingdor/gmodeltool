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

func (g *genSchemaCommander) GenTableSchema(ctx context.Context, tname string, name string, fields []*common.ColumnInfo) (code string, imports []common.Import, err error) {

	maxLen := g.maxColumnLen(fields)
	maxLen += 8
	structBuf := bytes.Buffer{}
	variable := bytes.Buffer{}
	refFields := bytes.Buffer{}
	var typeName = fmt.Sprintf("%sSchemaType", name)
	structBuf.WriteString(fmt.Sprintf("type %s struct{\n\t_alias string\n", typeName))
	variable.WriteString(fmt.Sprintf("var %sSchema *%s = &%s{\n", name, typeName, typeName))
	for _, column := range fields {
		field := column.Field
		fillName := utils.PadLeftRightSpaces(column.Name, 4, maxLen)
		structBuf.WriteString(fmt.Sprintf("\t// %s %s \n", column.Name, strings.ReplaceAll(column.Comment, "\n", "\n    //")))
		structBuf.WriteString(fmt.Sprintf("\t%s gmodel.Field\n", fillName))
		variable.WriteString(fmt.Sprintf("\t%s :gmodel.NewField(\"%s\", \"%s\", %v, %v),\n",
			fillName, field.Name(), field.Type(), field.IsNullable(), field.IsPK()))
		//&s.Id,
		refFields.WriteString(fmt.Sprintf("\t&s.%s,\n", column.Name))
	}
	structBuf.WriteString("\n}")
	variable.WriteString("\n}")

	if code, err = template.ReadFS("files/schema.go.template"); err == nil {
		code = strings.ReplaceAll(code, "{$schemaType}", structBuf.String())
		code = strings.ReplaceAll(code, "{$schema}", variable.String())
		code = strings.ReplaceAll(code, "{$schemaTypeName}", typeName)
		code = strings.ReplaceAll(code, "{$tableName}", fmt.Sprintf("%q", tname))
		code = strings.ReplaceAll(code, "{$refFields}", refFields.String())
	}
	imports = []common.Import{
		{Path: "github.com/lingdor/gmodel"},
		//{Path: "github.com/lingdor/gmodel/orm"},
		{Path: "fmt"},
	}
	return
}
