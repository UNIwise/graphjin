package sdata

import (
	"fmt"
	"strings"
)

// String returns a string representation of the DBTable
func (ti *DBTable) String() string {
	return ti.Schema + "." + ti.Name
}

// String returns a string representation of the DBColumn
func (col DBColumn) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s.%s.%s", col.Schema, col.Table, col.Name)
	fmt.Fprintf(&sb, " [id:%d, type:%v, array:%t, notNull:%t, fulltext:%t]",
		col.ID, col.Type, col.Array, col.NotNull, col.FullText)

	if col.FKeyCol != "" {
		fmt.Fprintf(&sb, " -> %s.%s.%s",
			col.FKeySchema, col.FKeyTable, col.FKeyCol)
	}
	return sb.String()
}

// String returns a string representation of the DBFunction
func (fn DBFunction) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s.%s", fn.Schema, fn.Name)
	fmt.Fprintf(&sb, " [type:%v, agg:%t] (", fn.Type, fn.Agg)

	for _, v := range fn.Inputs {
		if v.Name == "" {
			fmt.Fprintf(&sb, "%d: %v [array:%t]", v.ID, v.Type, v.Array)
		} else {
			fmt.Fprintf(&sb, "%s: %v [array:%t]", v.Name, v.Type, v.Array)
		}
	}

	sb.WriteString(") => ")

	for _, v := range fn.Outputs {
		if v.Name == "" {
			fmt.Fprintf(&sb, "%d: %v [array:%t]", v.ID, v.Type, v.Array)
		} else {
			fmt.Fprintf(&sb, "%s: %v [array:%t]", v.Name, v.Type, v.Array)
		}
	}

	return sb.String()
}

// String returns a string representation of the DBRel
func (re *DBRel) String() string {
	return fmt.Sprintf("'%s' --(%s)--> '%s'",
		re.Left.Col.String(),
		re.Type,
		re.Right.Col.String())

}
