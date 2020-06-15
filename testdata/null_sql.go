package testdata

import (
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type NullTest struct {
	Bool      *bool                `db:"bool" json:"bool,omitempty"`
	Bool2     sql.NullBool         `db:"bool2" json:"bool2,omitempty"`
	Time      *time.Time           `db:"time" json:"time,omitempty"`
	Time2     sql.NullTime         `db:"time2" json:"time2,omitempty"`
	Timestamp *timestamp.Timestamp `db:"timestamp" json:"timestamp,omitempty"`
	String    *string              `db:"string" json:"string,omitempty"`
	String2   sql.NullString       `db:"string2" json:"string2,omitempty"`
	Float32   *float32             `db:"float32" json:"float32,omitempty"`
	Float64   *float64             `db:"float64" json:"float64,omitempty"`
	Float642  sql.NullFloat64      `db:"float642" json:"float642,omitempty"`
	Int       *int                 `db:"int" json:"int,omitempty"`
	Int32     *int32               `db:"int32" json:"int32,omitempty"`
	Int322    sql.NullInt32        `db:"int322" json:"int322,omitempty"`
	Int64     *int64               `db:"int64" json:"int64,omitempty"`
	Int642    sql.NullInt64        `db:"int642" json:"int642,omitempty"`
	Uint      *uint                `db:"uint" json:"uint,omitempty"`
	Uint32    *uint32              `db:"uint32" json:"uint32,omitempty"`
	Uint64    *uint64              `db:"uint64" json:"uint64,omitempty"`
}

var NullQueryPG = `
select
CAST ( null AS bool ) as "bool",
CAST ( null AS bool ) as "bool2",
CAST ( null AS time ) as "time",
CAST ( null AS time ) as "time2",
CAST ( null AS time ) as "timestamp",
CAST ( null AS text ) as "string",
CAST ( null AS text ) as "string2",
CAST ( null AS real) as "float32",
CAST ( null AS double precision ) as "float64",
CAST ( null AS double precision ) as "float642",
CAST ( null AS integer ) as "int",
CAST ( null AS integer ) as "int32",
CAST ( null AS integer ) as "int322",
CAST ( null AS bigint ) as "int64",
CAST ( null AS bigint ) as "int642",
CAST ( null AS integer ) as "uint",
CAST ( null AS integer ) as "uint32",
CAST ( null AS bigint ) as "uint64"
`

var NotNullQueryPG = `
select
CAST ( 1 AS bool ) as "bool",
CAST ( 1 AS bool ) as "bool2",
CAST ( '2006-01-02T15:04:05Z07:00' AS timestamp ) as "time",
CAST ( '2006-01-02T15:04:05Z07:00' AS timestamp ) as "time2",
CAST ( '2006-01-02T15:04:05Z07:00' AS date ) as "timestamp", 
CAST ( 1 AS text ) as "string",
CAST ( 1 AS text ) as "string2",
CAST ( 1 AS real) as "float32",
CAST ( 1 AS double precision ) as "float64",
CAST ( 1 AS double precision ) as "float642",
CAST ( 1 AS integer ) as "int",
CAST ( 1 AS integer ) as "int32",
CAST ( 1 AS integer ) as "int322",
CAST ( 1 AS bigint ) as "int64",
CAST ( 1 AS bigint ) as "int642",
CAST ( 1 AS integer ) as "uint",
CAST ( 1 AS integer ) as "uint32",
CAST ( 1 AS bigint ) as "uint64"
`

var NullQueryMySql = `
select
nullbool as "bool",
nullbool as "bool2",
CAST( null AS DATETIME ) as "time",
CAST( null AS DATETIME ) as "time2",
CAST( null AS DATE ) as "timestamp", 
CAST( null AS CHAR ) as "string",
CAST( null AS CHAR ) as "string2",
CAST( null AS FLOAT(32)) as "float32",
CAST( null AS DECIMAL(64)) as "float64",
CAST( null AS DECIMAL(64)) as "float642",
CAST( null AS UNSIGNED ) as "int",
CAST( null AS UNSIGNED ) as "int32",
CAST( null AS UNSIGNED ) as "int322",
CAST( null AS UNSIGNED ) as "int64",
CAST( null AS UNSIGNED ) as "int642",
CAST( null AS UNSIGNED ) as "uint",
CAST( null AS UNSIGNED ) as "uint32",
CAST( null AS UNSIGNED ) as "uint64"
from nullbool
limit 1
`

var NotNullQueryMySQL = `
select
true as "bool",
true as "bool2",
CAST( '2006-01-02T15:04:05Z07:00' AS DATETIME ) as "time",
CAST( '2006-01-02T15:04:05Z07:00' AS DATETIME ) as "time2",
CAST( '2006-01-02T15:04:05Z07:00' AS DATE ) as "timestamp", 
CAST( 1 AS CHAR ) as "string",
CAST( 1 AS CHAR ) as "string2",
CAST( 1 AS float(32) ) as "float32",
CAST( 1 AS decimal(64) ) as "float64",
CAST( 1 AS decimal(64) ) as "float642",
CAST( 1 AS UNSIGNED ) as "int",
CAST( 1 AS UNSIGNED ) as "int32",
CAST( 1 AS UNSIGNED ) as "int322",
CAST( 1 AS UNSIGNED ) as "int64",
CAST( 1 AS UNSIGNED ) as "int642",
CAST( 1 AS UNSIGNED ) as "uint",
CAST( 1 AS UNSIGNED ) as "uint32",
CAST( 1 AS UNSIGNED ) as "uint64"
`
