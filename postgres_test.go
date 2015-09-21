package barkup

import (
	"strings"
	"testing"
)

func Test_Postgres_Export_Pass(t *testing.T) {
	p := Postgres{
		Host:     "localhost",
		Port:     "5432",
		DB:       "test",
		Username: "postgres",
	}

	PGDumpCmd = "true"

	result := p.Export()
	expect(t, result.Error, (*Error)(nil))
}

func Test_Postgres_Export_FailDump(t *testing.T) {
	p := Postgres{
		Host:     "localhost",
		Port:     "5432",
		DB:       "test",
		Username: "postgres",
	}

	PGDumpCmd = "false"

	result := p.Export()
	refute(t, result.Error, (*Error)(nil))
}

func Test_Postgres_optionsDump(t *testing.T) {
	p := Postgres{
		Host:     "localhost",
		Port:     "5432",
		DB:       "test",
		Username: "postgres",
		Options:  []string{"-W"},
	}

	options := strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-d"), true)
	expect(t, strings.Contains(options, "-h"), true)
	expect(t, strings.Contains(options, "-p"), true)
	expect(t, strings.Contains(options, "-U"), true)
	expect(t, strings.Contains(options, "-W"), true)

	p.Host = ""
	options = strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-h"), false)

	p.Port = ""
	options = strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-p"), false)

	p.DB = ""
	options = strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-d"), false)

	p.Username = ""
	options = strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-U"), false)

	p.Options = []string{}
	options = strings.Join(p.dumpOptions(), " ")
	expect(t, strings.Contains(options, "-W"), false)
}
