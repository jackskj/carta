package carta_test

import (
	"log"
	"testing"

	td "github.com/jackskj/carta/testdata"
)

func iTestCasting(t *testing.T) {
	req := td.EmptyRequest{}
	_, err := reflectClient.TypeCasting(ctx, &req)
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
}

func TestIncorrectType(t *testing.T) {
	// these are use in a dynamic sql statement
	var type_values = []string{
		"created_on", // sql returns datetime
		"1.1",        // sql returns float
		"1",          // sql returns int
		"true",       // sql returns bool
		"'1'",        // sql returns text
	}
	for _, sql_type := range type_values {
		req := td.TypeRequest{TypeValue: sql_type}
		_, err := reflectClient.IncorrectTypes(ctx, &req)
		if err != nil {
			log.Println(err)
		}
	}
}

func TestNullType(t *testing.T) {
	req := td.EmptyRequest{}
	resp, err := reflectClient.NullTypes(ctx, &req)
	if err != nil {
		log.Println(resp)
	}
}
