package mock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {

	args := m.Called(ctx, options)

	return args.Get(0).(pgx.Tx), args.Error(1)

}
func (m *MockDB) QueryRow(ctx context.Context, query string, _args ...interface{}) pgx.Row {
	var params []interface{}
	params = append(params, ctx, query)
	params = append(params, _args...)

	args := m.Called(params...)

	return args.Get(0).(pgx.Row)

}
func (m *MockDB) Query(ctx context.Context, query string, _args ...interface{}) (pgx.Rows, error) {
	var params []interface{}
	params = append(params, ctx, query)
	params = append(params, _args...)

	args := m.Called(params...)

	return args.Get(0).(pgx.Rows), args.Error(1)
}
func (m *MockDB) Exec(ctx context.Context, sql string, _args ...interface{}) (pgconn.CommandTag, error) {
	var params []interface{}
	params = append(params, ctx, sql)
	params = append(params, _args...)

	args := m.Called(params...)

	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}
func (m *MockDB) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	var params []interface{}
	params = append(params, ctx, b)

	args := m.Called(params...)

	return args.Get(0).(pgx.BatchResults)
}

type MockTx struct {
	MockDB
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)

	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)

	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	var params []interface{}
	params = append(params, dest...)

	args := m.Called(params...)

	return args.Error(0)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Close() {
	m.Called()
	return
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...any) error {
	var params []interface{}
	params = append(params, dest...)

	args := m.Called(params...)

	return args.Error(0)
}

func (m *MockRows) Values() ([]any, error) {
	args := m.Called()
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}
