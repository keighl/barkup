package barkup

import (
	"strings"
	"testing"
)

func Test_RethinkDB_Export_Pass(t *testing.T) {
	m := RethinkDB{
		Connection: "localhost:28015",
		Name:       "test",
	}

	RethinkCmd = "true"

	result := m.Export()
	expect(t, result.Error, (*Error)(nil))
}

func Test_RethinkDB_Export_FailDump(t *testing.T) {
	m := RethinkDB{
		Connection: "localhost:28015",
		Name:       "test",
	}

	RethinkCmd = "false"

	result := m.Export()
	refute(t, result.Error, (*Error)(nil))
}

func Test_RethinkDB_optionsDump(t *testing.T) {
	m := RethinkDB{
		Connection: "localhost:28015",
		Name:       "test",
		Targets:    []string{"cheese", "milk"},
		Options:    []string{"-aAUTHKEY"},
	}

	optionsR := m.dumpOptions()
	expect(t, optionsR[0], "dump")
	options := strings.Join(optionsR, " ")
	expect(t, strings.Contains(options, "-clocalhost:28015"), true)
	expect(t, strings.Contains(options, "-aAUTHKEY"), true)
	expect(t, strings.Contains(options, "-echeese"), true)
	expect(t, strings.Contains(options, "-emilk"), true)
}
