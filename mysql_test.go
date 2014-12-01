package barkup

import (
  "testing"
  "errors"
)

func Test_MySQL_Export_Pass(t *testing.T) {
  m := MySQL{
    Host: "localhost",
    Port: "3306",
    DB: "test",
    User: "root",
    Password: "",
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
    Password: "",
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
    Password: "",
  }

  mysqlDump = func (x MySQL, path string) error { return nil }
  mysqlTar = func (x MySQL, path string, destPath string) error { return errors.New("***") }

  result := m.Export()
  refute(t, result.Error, nil)
}

