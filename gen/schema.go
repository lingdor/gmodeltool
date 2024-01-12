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

func (g *genSchemaCommander) GenSchema(ctx context.Context, w *os.File, table, objName string, fields []*common.ColumnInfo, start, end string) {

	if code, err := g.GenTableSchema(ctx, table, objName, fields); err == nil {

		w.WriteString(fmt.Sprintf("%s:%s\n", start, common.MD5(code)))
		w.WriteString(code)
		w.WriteString(fmt.Sprintf("\n%s", end))
	}
}

func (g *genSchemaCommander) GenTableSchema(ctx context.Context, tname string, name string, fields []*common.ColumnInfo) (code string, err error) {

	structBuf := bytes.Buffer{}
	variable := bytes.Buffer{}
	var typeName = fmt.Sprintf("%sSchemaType", name)
	structBuf.WriteString(fmt.Sprintf("type %s struct{\n", typeName))
	variable.WriteString(fmt.Sprintf("var %sSchema %s{\n", name, typeName))
	for _, column := range fields {
		field := column.Field

		structBuf.WriteString(fmt.Sprintf("    // %s %s \n", column.Name, column.Comment))
		structBuf.WriteString(fmt.Sprintf("    %s	gmodel.Field\n", column.Name))
		variable.WriteString(fmt.Sprintf("	%s:		gmodel.NewField(\"%s\",\"%s\",%v,%v),\n",
			column.Name, field.Name(), field.Type(), field.IsNullable(), field.IsPK()))
	}
	structBuf.WriteString("\n}")
	variable.WriteString("\n}")

	if code, err = template.ReadFS("files/schema.go.template"); err == nil {
		code = strings.ReplaceAll(code, "{$schemaType}", structBuf.String())
		code = strings.ReplaceAll(code, "{$schema}", variable.String())
		code = strings.ReplaceAll(code, "{$schemaTypeName}", typeName)
		code = strings.ReplaceAll(code, "{$tableName}", fmt.Sprintf("%q", tname))
	}
	return
}
