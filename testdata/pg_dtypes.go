package testdata

var PGDtypes = `
bigint                       --- "~18 digit integer, 8-byte storage"
bit                          --- fixed-length bit string
bit varying                  --- variable-length bit string
boolean                      --- "boolean, 'true'/'false'"
character                    --- "char(length), blank-padded string, fixed storage length"
character varying            --- "varchar(length), non-blank-padded string, variable storage length"
cid                          --- "command identifier type, sequence in transaction id"
cidr                         --- "network IP address/netmask, network address"
date                         --- date
daterange                    --- range of dates
double precision             --- "double-precision floating point number, 8-byte storage"
inet                         --- "IP address/netmask, host address, netmask optional"
int4range                    --- range of integers
int8range                    --- range of bigints
integer                      --- "-2 billion to 2 billion integer, 4-byte storage"
numeric                      --- "numeric(precision, decimal), arbitrary precision number"
numrange                     --- range of numerics
oid                          --- "object identifier(oid), maximum 4 billion"
real                         --- "single-precision floating point number, 4-byte storage"
smallint                     --- "-32 thousand to 32 thousand, 2-byte storage"
text                         --- "variable-length string, no limit specified"
timestamp without time zone  --- date and time
timestamp with time zone     --- date and time with time zone
time without time zone       --- time of day
time with time zone          --- time of day with time zone
tstzrange                    --- range of timestamps with time zone
uuid                         --- UUID datatype
xml                          --- XML content
`
