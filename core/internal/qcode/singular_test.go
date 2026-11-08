package qcode

import (
	"testing"

	"github.com/dosco/graphjin/core/v3/internal/sdata"
)

func eqExp(col sdata.DBColumn) *Exp {
	ex := newExpOp(OpEquals)
	ex.Left.Col = col

	return ex
}

func andExp(children ...*Exp) *Exp {
	ex := newExpOp(OpAnd)
	ex.Children = children

	return ex
}

func orExp(children ...*Exp) *Exp {
	ex := newExpOp(OpOr)
	ex.Children = children

	return ex
}

var (
	colID = sdata.DBColumn{Name: "id", PrimaryKey: true, UniqueKey: true}
	colA  = sdata.DBColumn{Name: "a"}
	colB  = sdata.DBColumn{Name: "b"}
	colC  = sdata.DBColumn{Name: "c"}
)

// table with a primary key on id and a composite unique key on (a, b).
func indexedTable() sdata.DBTable {
	return sdata.DBTable{
		Name:          "t",
		Columns:       []sdata.DBColumn{colID, colA, colB, colC},
		UniqueIndices: [][]string{{"id"}, {"a", "b"}},
	}
}

func TestSingularByConstraints(t *testing.T) {
	cases := []struct {
		name   string
		ti     sdata.DBTable
		rel    sdata.DBRel
		where  *Exp
		paging Paging
		want   bool
	}{
		{name: "limit 1 alone", ti: indexedTable(), paging: Paging{Limit: 1}, want: true},
		{name: "limit 2 alone", ti: indexedTable(), paging: Paging{Limit: 2}, want: false},
		{name: "first 1 (forward cursor pagination) is not a limit 1", ti: indexedTable(), paging: Paging{Type: PTForward, Limit: 1}, want: false},
		{name: "last 1 (backward cursor pagination) is not a limit 1", ti: indexedTable(), paging: Paging{Type: PTBackward, Limit: 1}, want: false},
		{name: "limit 1 with a cursor is not singular", ti: indexedTable(), paging: Paging{Limit: 1, Cursor: true}, want: false},
		{name: "no where and no limit", ti: indexedTable(), want: false},
		{name: "primary key equality", ti: indexedTable(), where: eqExp(colID), want: true},
		{name: "non-unique column equality", ti: indexedTable(), where: eqExp(colC), want: false},
		{name: "composite key, both columns", ti: indexedTable(), where: andExp(eqExp(colA), eqExp(colB)), want: true},
		{name: "composite key, one column", ti: indexedTable(), where: eqExp(colA), want: false},
		{name: "composite key found through nested and", ti: indexedTable(),
			where: andExp(eqExp(colC), andExp(eqExp(colA), eqExp(colB))), want: true},
		{name: "composite key spread over or", ti: indexedTable(), where: orExp(eqExp(colA), eqExp(colB)), want: false},
		{name: "unique column under or does not count", ti: indexedTable(), where: orExp(eqExp(colID), eqExp(colC)), want: false},
		{
			name:  "join column completes a composite key",
			ti:    indexedTable(),
			rel:   sdata.DBRel{Left: sdata.DBRelLeft{Col: colB}},
			where: eqExp(colA),
			want:  true,
		},
		{
			// The column flags mark a member of a composite key as unique, so
			// they must not make a selection singular.
			name:  "no index metadata: a flagged column does not make it singular",
			ti:    sdata.DBTable{Name: "t", Columns: []sdata.DBColumn{colID, colC}},
			where: eqExp(colID),
			want:  false,
		},
		{
			name:  "no index metadata: limit 1 still applies",
			ti:    sdata.DBTable{Name: "t", Columns: []sdata.DBColumn{colID, colC}},
			where: eqExp(colC), paging: Paging{Limit: 1},
			want: true,
		},
		{name: "non-equality operator is ignored", ti: indexedTable(), where: func() *Exp {
			ex := newExpOp(OpGreaterThan)
			ex.Left.Col = colID

			return ex
		}(), want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isSingularByConstraints(c.ti, c.rel, c.where, c.paging); got != c.want {
				t.Errorf("isSingularByConstraints = %v, want %v", got, c.want)
			}
		})
	}
}

func TestSingularByConstraintsAppliesToMySQLFamilyOnly(t *testing.T) {
	for _, dbType := range []string{"mysql", "mariadb"} {
		if !SingularByConstraintsApplies(dbType) {
			t.Errorf("%s must use the constraint-based singular rule", dbType)
		}
	}

	// Every other dialect keeps upstream's behaviour exactly.
	for _, dbType := range []string{"postgres", "", "sqlite", "mssql", "redshift", "oracle", "snowflake", "bigquery"} {
		if SingularByConstraintsApplies(dbType) {
			t.Errorf("%q must keep upstream behaviour", dbType)
		}
	}
}
