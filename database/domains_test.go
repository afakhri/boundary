package database

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v4"
)

func TestPublicIdDomain(t *testing.T) {
	const (
		createTable = `
create table test_table (
  id bigint generated always as identity primary key,
  public_id icu_public_id
);
`
		insert = `
insert into test_table (public_id)
values ($1)
returning id;
`
	)

	cleanup, connURL := PreparePostgresTestContainer(t)
	defer cleanup()

	db, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("query: \n%s\n error: %s", createTable, err)
	}

	failTests := []string{
		" ",
		"bar",
		"0000000001000000000200000000031",
	}
	for _, tt := range failTests {
		value := tt
		t.Run(tt, func(t *testing.T) {
			t.Logf("insert value: %q", value)
			if _, err := db.Query(insert, value); err == nil {
				t.Errorf("want error, got no error for inserting public_id: %s", value)
			}
		})
	}

	okTests := []string{
		"l1Ocw0TpHn800CekIxIXlmQqRDgFDfYl",
		"00000000010000000002000000000312",
		"00000000010000000002000000000032",
		"12345678901234567890123456789012",
		"ec2_12345678901234567890123456789012",
		"prj_12345678901234567890123456789012",
	}
	for _, tt := range okTests {
		value := tt
		t.Run(tt, func(t *testing.T) {
			t.Logf("insert value: %q", value)
			if _, err := db.Query(insert, value); err != nil {
				t.Errorf("%s: want no error, got error %v", value, err)
			}
		})
	}
}