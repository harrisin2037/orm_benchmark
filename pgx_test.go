package orm_benchmark

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
)

type Pgx struct {
}

func PgxTestBenchORMInstance() BenchORMInstance {
	return &Pgx{}
}

func (ins *Pgx) ConnTest(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn := ins.openConn(b)
			ins.closeConn(b, conn)
		}
	})
}

func (ins *Pgx) openConnDB(tb testing.TB) *sql.DB {

	config, err := pgx.ParseConfig(os.Getenv("ORM_TEST_DATABASE"))
	require.Nil(tb, err)

	config.BuildStatementCache = nil

	conf := *config

	return stdlib.OpenDB(conf)
}

func (ins *Pgx) closeConnDB(tb testing.TB, db *sql.DB) {
	err := db.Close()
	require.NoError(tb, err)
}

func (ins *Pgx) openConn(tb testing.TB) *pgx.Conn {

	config, err := pgx.ParseConfig(os.Getenv("ORM_TEST_DATABASE"))
	require.Nil(tb, err)

	config.BuildStatementCache = nil

	conn, err := pgx.ConnectConfig(context.Background(), config)
	require.Nil(tb, err, fmt.Sprintf("ORM_TEST_DATABASE: %s", os.Getenv("ORM_TEST_DATABASE")))

	return conn
}

func (ins *Pgx) closeConn(tb testing.TB, conn *pgx.Conn) {
	err := conn.Close(context.Background())
	require.Nil(tb, err)
}
