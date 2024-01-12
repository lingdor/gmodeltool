package common

import "github.com/lingdor/gmodel"

type ColumnInfo struct {
	Field   *gmodel.FieldInfo
	Comment string
	Name    string
}
