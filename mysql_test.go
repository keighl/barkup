package barkup

import (
  "testing"
  "errors"
  "strings"
)

func Test_MySQL_Export_Pass(t *testing.T) {
  m := MySQL{
    Host: "localhost",
    Port: "3306",
    DB: "test",
    User: "root",
    Password: "cheese",
  }

  mysqlDump = func (x MySQL, path string) error { return nil }
  mysqlTar = func (x MySQL, path string, destPath string) error { return nil }

  result := m.Export()
  expect(t, result.Error, nil)
}

func Test_MySQL_Export_FailDump(t *testing.T) {
  m := MySQL{
    Host: "localhost",
    Port: "3306",
    DB: "test",
    User: "root",
    Password: "cheese",
  }

  mysqlDump = func (x MySQL, path string) error { return errors.New("***") }
  mysqlTar = func (x MySQL, path string, destPath string) error { return nil }

  result := m.Export()
  refute(t, result.Error, nil)
}

func Test_MySQL_Export_FailTar(t *testing.T) {
  m := MySQL{
    Host: "localhost",
    Port: "3306",
    DB: "test",
    User: "root",
    Password: "cheese",
  }

  mysqlDump = func (x MySQL, path string) error { return nil }
  mysqlTar = func (x MySQL, path string, destPath string) error { return errors.New("***") }

  result := m.Export()
  refute(t, result.Error, nil)
}

func Test_MySQL_optionsDump(t *testing.T) {
  m := MySQL{
    Host: "localhost",
    Port: "3306",
    DB: "test",
    User: "root",
    Password: "cheese",
    Options: []string{"--skip-extended-insert"},
  }

  options := strings.Join(m.dumpOptions(), " ")
  expect(t, strings.Contains(options, "-h"), true)
  expect(t, strings.Contains(options, "-P"), true)
  expect(t, strings.Contains(options, "-u"), true)
  expect(t, strings.Contains(options, "-p"), true)
  expect(t, strings.Contains(options, "--skip-extended-insert"), true)
  expect(t, strings.Contains(options, m.DB), true)
}

