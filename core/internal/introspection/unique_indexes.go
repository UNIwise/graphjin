package introspection

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// mysqlUniqueIndexesStmt lists the column sets of every unique and primary key
// index, in index order. NON_UNIQUE = 0 covers both PRIMARY and UNIQUE keys.
const mysqlUniqueIndexesStmt = `
SELECT TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, COLUMN_NAME
FROM information_schema.STATISTICS
WHERE NON_UNIQUE = 0
  AND TABLE_SCHEMA NOT IN ('_graphjin', 'information_schema', 'performance_schema', 'mysql', 'sys')
ORDER BY TABLE_SCHEMA, TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`

// discoverUniqueIndexes returns, per "schema:table", the column set of each
// unique or primary key index. It exists so singular detection can recognise
// composite keys, which the per-column UniqueKey flag cannot express.
func discoverUniqueIndexes(ctx context.Context, db *sql.DB) (map[string][][]string, error) {
	qctx, cancel := context.WithTimeout(ctx, introspectionQueryTimeout)
	defer cancel()

	start := time.Now()
	rows, err := db.QueryContext(qctx, mysqlUniqueIndexesStmt)
	recordDiscoveryQuery(ctx, "unique_indexes", mysqlUniqueIndexesStmt, time.Since(start))

	if err != nil {
		return nil, fmt.Errorf("error fetching unique indexes: %w", err)
	}
	defer rows.Close()

	type indexKey struct{ schema, table, index string }

	var (
		order   []indexKey
		columns = map[indexKey][]string{}
	)

	for rows.Next() {
		var schema, table, index, column string
		if err := rows.Scan(&schema, &table, &index, &column); err != nil {
			return nil, fmt.Errorf("error scanning unique indexes: %w", err)
		}

		k := indexKey{schema, table, index}
		if _, seen := columns[k]; !seen {
			order = append(order, k)
		}

		columns[k] = append(columns[k], column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading unique indexes: %w", err)
	}

	out := make(map[string][][]string)

	for _, k := range order {
		tableKey := k.schema + ":" + k.table
		out[tableKey] = append(out[tableKey], columns[k])
	}

	return out, nil
}
