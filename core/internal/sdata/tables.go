package sdata

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/sync/errgroup"
)

type DBInfo struct {
	Type    string
	Version int
	Schema  string
	Name    string

	Tables    []DBTable      `hash:"set"`
	Functions []DBFunction   `hash:"set"`
	VTables   []VirtualTable `hash:"set"`
	colMap    map[string]int `hash:"-"`
	tableMap  map[string]int `hash:"-"`
	hash      uint64         `hash:"-"`
}

type DBIndices map[string][]DBColumnIndex

type DBTable struct {
	Schema       string
	Name         string
	Type         string
	Columns      []DBColumn `hash:"set"`
	PrimaryCol   DBColumn
	SecondaryCol DBColumn
	FullText     []DBColumn `hash:"set"`
	Blocked      bool
	colMap       map[string]int `hash:"-"`
	IndexColumns DBIndices
	Indices      DBIndices
}

type DBIndexTable struct {
	Columns DBIndices
	Indices DBIndices
}

type DBColumnIndex struct {
	Schema     string
	Constraint string
	Table      string
	Type       string
	Column     string
	Composite  bool
}

type VirtualTable struct {
	Name       string
	IDColumn   string
	TypeColumn string
	FKeyColumn string
}

type st struct {
	schema, table string
}

func GetDBInfo(
	db *sql.DB,
	dbType string,
	blockList []string) (*DBInfo, error) {

	var dbVersion int
	var dbSchema, dbName string
	var cols []DBColumn
	var funcs []DBFunction
	var tableIndices map[string]DBIndexTable
	var err error

	g := errgroup.Group{}

	g.Go(func() error {
		var row *sql.Row
		switch dbType {
		case "mysql":
			row = db.QueryRow(mysqlInfo)
		default:
			row = db.QueryRow(postgresInfo)
		}

		if err := row.Scan(&dbVersion, &dbSchema, &dbName); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		if cols, err = DiscoverColumns(db, dbType, blockList); err != nil {
			return err
		}

		if funcs, err = DiscoverFunctions(db, blockList); err != nil {
			return err
		}

		if tableIndices, err = DiscoverIndices(db); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	di := NewDBInfo(
		dbType,
		dbVersion,
		dbSchema,
		dbName,
		cols,
		funcs,
		tableIndices,
		blockList,
	)

	di.hash, err = hashstructure.Hash(di, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}

	return di, nil
}

func NewDBInfo(
	dbType string,
	dbVersion int,
	dbSchema string,
	dbName string,
	cols []DBColumn,
	funcs []DBFunction,
	tableIndices map[string]DBIndexTable,
	blockList []string) *DBInfo {

	di := &DBInfo{
		Type:      dbType,
		Version:   dbVersion,
		Schema:    dbSchema,
		Name:      dbName,
		Functions: funcs,
		colMap:    make(map[string]int),
		tableMap:  make(map[string]int),
	}

	tm := make(map[st][]DBColumn)
	for i := range cols {
		c := cols[i]
		di.colMap[(c.Schema + ":" + c.Table + ":" + c.Name)] = i

		k1 := st{c.Schema, c.Table}
		tm[k1] = append(tm[k1], c)
	}

	for k, tcols := range tm {
		ti := NewDBTable(k.schema, k.table, "", tcols, tableIndices[k.table].Columns, tableIndices[k.table].Indices)
		if strings.HasPrefix(ti.Name, "_gj_") {
			continue
		}
		ti.Blocked = isInList(ti.Name, blockList)
		di.AddTable(ti)
	}

	return di
}

func NewDBTable(schema, name, _type string, cols []DBColumn, columnIndices map[string][]DBColumnIndex, indices map[string][]DBColumnIndex) DBTable {
	ti := DBTable{
		Schema:       schema,
		Name:         name,
		Type:         _type,
		Columns:      cols,
		colMap:       make(map[string]int, len(cols)),
		IndexColumns: columnIndices,
		Indices:      indices,
	}

	for i, c := range cols {
		switch {
		case c.FullText:
			ti.FullText = append(ti.FullText, c)

		case c.PrimaryKey:
			ti.PrimaryCol = c

		}
		ti.colMap[c.Name] = i
	}
	return ti
}

func (di *DBInfo) AddTable(t DBTable) {
	for i, c := range t.Columns {
		di.colMap[(c.Schema + ":" + c.Table + ":" + c.Name)] = i
	}

	i := len(di.Tables)
	di.Tables = append(di.Tables, t)
	di.tableMap[(t.Schema + ":" + t.Name)] = i
}

func (di *DBInfo) GetColumn(schema, table, column string) (*DBColumn, error) {
	t, err := di.GetTable(schema, table)
	if err != nil {
		return nil, err
	}

	cid, ok := t.colMap[column]
	if !ok {
		return nil, fmt.Errorf("column: '%s.%s.%s' not found", schema, table, column)
	}

	return &t.Columns[cid], nil
}

func (di *DBInfo) GetTable(schema, table string) (*DBTable, error) {
	tid, ok := di.tableMap[(schema + ":" + table)]
	if !ok {
		return nil, fmt.Errorf("table: '%s.%s' not found", schema, table)
	}

	return &di.Tables[tid], nil
}

type DBColumn struct {
	ID         int32
	Name       string
	Type       string
	Array      bool
	NotNull    bool
	PrimaryKey bool
	UniqueKey  bool
	FullText   bool
	FKeySchema string
	FKeyTable  string
	FKeyCol    string
	Blocked    bool
	Table      string
	Schema     string
}

func DiscoverColumns(db *sql.DB, dbtype string, blockList []string) ([]DBColumn, error) {
	var sqlStmt string

	switch dbtype {
	case "mysql":
		sqlStmt = mysqlColumnsStmt
	default:
		sqlStmt = postgresColumnsStmt
	}

	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error fetching columns: %s", err)
	}
	defer rows.Close()

	cmap := make(map[string]DBColumn)

	for rows.Next() {
		var c DBColumn

		err = rows.Scan(&c.Schema, &c.Table, &c.Name, &c.Type, &c.NotNull, &c.PrimaryKey, &c.UniqueKey, &c.Array, &c.FullText, &c.FKeySchema, &c.FKeyTable, &c.FKeyCol)

		if err != nil {
			return nil, err
		}

		k := (c.Schema + ":" + c.Table + ":" + c.Name)
		v, ok := cmap[k]
		if !ok {
			v = c
			v.ID = int32(len(cmap))
			if strings.HasPrefix(v.Table, "_gj_") {
				continue
			}
			v.Blocked = isInList(v.Name, blockList)
		}
		if c.Type != "" {
			v.Type = c.Type
		}
		if c.PrimaryKey {
			v.PrimaryKey = true
			v.UniqueKey = true
		}
		if c.NotNull {
			v.NotNull = true
		}
		if c.UniqueKey {
			v.UniqueKey = true
		}
		if c.Array {
			v.Array = true
		}
		if c.FullText {
			v.FullText = true
		}
		if c.FKeySchema != "" {
			v.FKeySchema = c.FKeySchema
		}
		if c.FKeyTable != "" {
			v.FKeyTable = c.FKeyTable
		}
		if c.FKeyCol != "" {
			v.FKeyCol = c.FKeyCol
		}
		cmap[k] = v
	}

	var cols []DBColumn
	for _, c := range cmap {
		cols = append(cols, c)
	}

	return cols, nil
}

type DBFunction struct {
	Name   string
	Params []DBFuncParam
}

type DBFuncParam struct {
	ID   int
	Name sql.NullString
	Type string
}

func DiscoverIndices(db *sql.DB) (map[string]DBIndexTable, error) {
	rows, err := db.Query(mysqlIndexInfoStmt)
	if err != nil {
		return nil, fmt.Errorf("error fetching index info: %s", err)
	}
	defer rows.Close()

	dbIndexTables := make(map[string]DBIndexTable)

	for rows.Next() {
		var ci DBColumnIndex
		ci.Composite = false
		err = rows.Scan(&ci.Schema, &ci.Constraint, &ci.Table, &ci.Type, &ci.Column)
		if err != nil {
			return nil, err
		}

		if _, ok := dbIndexTables[ci.Table]; !ok {
			dbIndexTables[ci.Table] = DBIndexTable{
				Columns: make(DBIndices),
				Indices: make(DBIndices),
			}
		}

		dbIndexTables[ci.Table].Indices[ci.Constraint] = append(dbIndexTables[ci.Table].Indices[ci.Constraint], ci)
		dbIndexTables[ci.Table].Columns[ci.Column] = append(dbIndexTables[ci.Table].Columns[ci.Column], ci)
	}

	return dbIndexTables, nil
}

func DiscoverFunctions(db *sql.DB, blockList []string) ([]DBFunction, error) {
	rows, err := db.Query(functionsStmt)
	if err != nil {
		return nil, fmt.Errorf("Error fetching functions: %s", err)
	}
	defer rows.Close()

	var funcs []DBFunction
	fm := make(map[string]int)

	parameterIndex := 1
	for rows.Next() {
		var fn, fid string
		fp := DBFuncParam{}

		err = rows.Scan(&fn, &fid, &fp.Type, &fp.Name, &fp.ID)
		if err != nil {
			return nil, err
		}

		if !fp.Name.Valid {
			fp.Name.String = fmt.Sprintf("%d", parameterIndex)
			fp.Name.Valid = true
		}

		if i, ok := fm[fid]; ok {
			funcs[i].Params = append(funcs[i].Params, fp)
		} else {
			if isInList(fn, blockList) {
				continue
			}
			funcs = append(funcs, DBFunction{Name: fn, Params: []DBFuncParam{fp}})
			fm[fid] = len(funcs) - 1
		}
		parameterIndex++
	}

	return funcs, nil
}

func (di *DBInfo) Hash() uint64 {
	return di.hash
}

func isInList(val string, s []string) bool {
	for _, v := range s {
		regex := fmt.Sprintf("^%s$", v)
		if matched, _ := regexp.MatchString(regex, val); matched {
			return true
		}
	}
	return false
}
