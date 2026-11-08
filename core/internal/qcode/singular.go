package qcode

import "github.com/dosco/graphjin/core/v3/internal/sdata"

// isSingularByConstraints reports whether a selection can return at most one
// row, so its result is an object rather than an array. That is the case when:
//
//   - the selection has an explicit `limit: 1` (not cursor pagination), or
//   - the columns it matches exactly (equality in the where clause, including
//     role filters, joined by AND, plus the column it is joined to its parent
//     on) cover every column of a unique or primary key index.
//
// Index metadata comes from UniqueIndices only. The per-column UniqueKey and
// PrimaryKey flags are not used: they also mark a single member of a composite
// key as unique. A table without index metadata is therefore never singular by
// constraints, which is the safe direction. Schema sources that carry no index
// metadata (the DDL cache) are not used for MySQL/MariaDB, see
// core/init_multidb.go.
func isSingularByConstraints(ti sdata.DBTable, rel sdata.DBRel, where *Exp, paging Paging) bool {
	// An explicit `limit: 1`. Cursor pagination (`first: 1`, `last: 1`) also has a
	// limit of 1 but returns pages of rows, so it stays an array.
	if paging.Limit == 1 && paging.Type == PTOffset && !paging.Cursor {
		return true
	}

	matched := make(map[string]sdata.DBColumn)

	if where != nil {
		collectEqualityColumns(where, matched)
	}

	if name := rel.Left.Col.Name; name != "" {
		matched[name] = rel.Left.Col
	}

	if len(matched) == 0 {
		return false
	}

	for _, index := range ti.UniqueIndices {
		if len(index) != 0 && coversAll(matched, index) {
			return true
		}
	}

	return false
}

func coversAll(matched map[string]sdata.DBColumn, index []string) bool {
	for _, name := range index {
		if _, ok := matched[name]; !ok {
			return false
		}
	}

	return true
}

// collectEqualityColumns gathers the columns matched with `=` through AND
// chains. Anything under OR, or any other operator, does not pin a column.
func collectEqualityColumns(ex *Exp, out map[string]sdata.DBColumn) {
	switch ex.Op {
	case OpEquals:
		if name := ex.Left.Col.Name; name != "" {
			out[name] = ex.Left.Col
		}
	case OpAnd:
		for _, child := range ex.Children {
			collectEqualityColumns(child, out)
		}
	}
}

// SingularByConstraintsApplies limits the rule to MySQL and MariaDB. Every
// other dialect keeps upstream's behaviour, so upstream's tests for them stay
// valid.
func SingularByConstraintsApplies(dbType string) bool {
	return dbType == "mysql" || dbType == "mariadb"
}
