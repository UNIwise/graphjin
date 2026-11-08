package sdata

import (
	"reflect"
	"testing"
)

func TestDBInfoSnapshotRoundTripKeepsUniqueIndices(t *testing.T) {
	columns := []DBColumn{
		{ID: 1, Name: "flowId", Type: "int", Schema: "wf", Table: "manager", PrimaryKey: false},
		{ID: 2, Name: "userId", Type: "int", Schema: "wf", Table: "manager"},
	}

	original := NewDBInfo("mysql", 80023, "wf", "wf", columns, nil, nil)
	original.Tables[0].UniqueIndices = [][]string{{"id"}, {"flowId", "userId"}}

	data, err := MarshalDBInfoSnapshot(original)
	if err != nil {
		t.Fatal(err)
	}

	restored, err := UnmarshalDBInfoSnapshot(data)
	if err != nil {
		t.Fatal(err)
	}

	table, err := restored.GetTable("wf", "manager")
	if err != nil {
		t.Fatal(err)
	}

	want := [][]string{{"id"}, {"flowId", "userId"}}
	if !reflect.DeepEqual(table.UniqueIndices, want) {
		t.Errorf("UniqueIndices after a snapshot round trip = %v, want %v", table.UniqueIndices, want)
	}
}

// A snapshot written before UniqueIndices existed (format version 1) would load
// without them and silently change which selections are singular, so it must be
// rejected and discovery run again.
func TestDBInfoSnapshotWithoutUniqueIndicesIsRejected(t *testing.T) {
	if _, err := UnmarshalDBInfoSnapshot([]byte(`{"version":1,"type":"mysql","tables":[]}`)); err == nil {
		t.Fatal("a version 1 snapshot lacks unique index metadata and must be rejected")
	}
}
