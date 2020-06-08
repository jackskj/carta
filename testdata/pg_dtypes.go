package testdata

import (
	"time"
)

var PGDTypesQuery = `
select
cast ( 1                                       as  bigint )                       as  "bigint",                      --- "~18 digit integer, 8-byte storage"
cast ( 1                                       as  bit )                          as  "bit",                         --- fixed-length bit string
cast ( 1                                       as  boolean )                      as  "boolean",                     --- "boolean, 'true'/'false'"
cast ( 'a'                                     as  character )                    as  "character",                   --- "char(length ) , blank-padded string, fixed storage length"
cast ( 'a'                                     as  character varying )            as  "character_varying",           --- "varchar(length ) , non-blank-padded string, variable storage length"
cast ( '1.1.1.1'                               as  cidr )                         as  "cidr",                        --- "network IP address/netmask, network address"
cast ( '2004-10-19'                            as  date )                         as  "date",                        --- date
cast (1                                        as  double precision )             as  "double_precision",            --- "double-precision floating point number, 8-byte storage"
cast ( 1                                       as  integer )                      as  "integer",                     --- "-2 billion to 2 billion integer, 4-byte storage"
cast ( 1                                       as  numeric )                      as  "numeric",                     --- "numeric(precision, decimal ) , arbitrary precision number"
cast ( 1                                       as  oid )                          as  "oid",                         --- "object identifier(oid ) , maximum 4 billion"
cast ( 1                                       as  real )                         as  "real",                        --- "single-precision floating point number, 4-byte storage"
cast ( 1                                       as  smallint )                     as  "smallint",                    --- "-32 thousand to 32 thousand, 2-byte storage"
cast ( 'a'                                     as  text )                         as  "text",                        --- "variable-length string, no limit specified"
cast (TIMESTAMP '2004-10-19 10:23:54'          as  timestamp without time zone )  as  "timestamp_without_time_zone", --- date and time
cast ( TIMESTAMP '2004-10-19 10:23:54+02'      as  timestamp with time zone )     as  "timestamp_with_time_zone",    --- date and time with time zone
cast ( '04:05:06'                              as  time without time zone )       as  "time_without_time_zone",      --- time of day
cast ( '04:05:06 PST'                          as  time with time zone )          as  "time_with_time_zone",         --- time of day with time zone
cast ( 'A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11'  as  uuid )                         as  "uuid",                        --- UUID datatype
cast ( 'a'                                     as  xml )                          as  "xml"                          --- XML content
--- cast ( 1                                       as  bit varying )                  as  "bit_varying",                 --- variable-length bit string
--- cast (                                         as  tstzrange )                    as  "tstzrange",                   --- range of timestamps with time zone
--- cast (                                         as  numrange )                     as  "numrange",                    --- range of numerics
--- cast (                                         as  inet )                         as  "inet",                        --- "IP address/netmask, host address, netmask optional"
--- cast (                                         as  int4range )                    as  "int4range",                   --- range of integers
--- cast (                                         as  int8range )                    as  "int8range",                   --- range of bigints
--- cast ( 1                                       as  cid )                          as  "cid",                         --- "command identifier type, sequence in transaction id"
--- cast (                                         as  daterange )                    as  "daterange",                   --- range of dates
`

type PGDTypes struct {
	Bigint                   int       `db:"bigint" json:"bigint,omitempty"`
	Bit                      string    `db:"bit" json:"bit,omitempty"`
	Boolean                  bool      `db:"boolean" json:"boolean,omitempty"`
	Character                string    `db:"character" json:"character,omitempty"`
	CharacterVarying         string    `db:"character_varying" json:"character_varying,omitempty"`
	Cidr                     string    `db:"cidr" json:"cidr,omitempty"`
	Date                     time.Time `db:"date" json:"date,omitempty"`
	DoublePrecision          float64   `db:"double_precision" json:"double_precision,omitempty"`
	Integer                  int       `db:"integer" json:"integer,omitempty"`
	Numeric                  int       `db:"numeric" json:"numeric,omitempty"`
	Oid                      int       `db:"oid" json:"oid,omitempty"`
	Real                     float32   `db:"real" json:"real,omitempty"`
	Smallint                 int       `db:"smallint" json:"smallint,omitempty"`
	Text                     string    `db:"text" json:"text,omitempty"`
	TimestampWithoutTimeZone time.Time `db:"timestamp_without_time_zone" json:"timestamp_without_time_zone,omitempty"`
	TimestampWithTimeZone    time.Time `db:"timestamp_with_time_zone" json:"timestamp_with_time_zone,omitempty"`
	TimeWithoutTimeZone      time.Time `db:"time_without_time_zone" json:"time_without_time_zone,omitempty"`
	TimeWithTimeZone         time.Time `db:"time_with_time_zone" json:"time_with_time_zone,omitempty"`
	Uuid                     string    `db:"uuid" json:"uuid,omitempty"`
	Xml                      string    `db:"xml" json:"xml,omitempty"`
}
