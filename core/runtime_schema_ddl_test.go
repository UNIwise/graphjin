package core

import "testing"

// The DDL cache has no unique-index metadata, so MySQL and MariaDB must never
// be initialised from it; every other dialect keeps upstream's behaviour.
func TestRuntimeSchemaDDLUsable(t *testing.T) {
	for _, dbType := range []string{"mysql", "mariadb"} {
		if runtimeSchemaDDLUsable(dbType) {
			t.Errorf("%s must not be initialised from the DDL cache", dbType)
		}
	}

	for _, dbType := range []string{"postgres", "", "sqlite", "mssql"} {
		if !runtimeSchemaDDLUsable(dbType) {
			t.Errorf("%q must keep using the DDL cache", dbType)
		}
	}
}
