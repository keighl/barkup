package barkup

import (
  "fmt"
  "time"
  "os/exec"
  "os"
)

// MySQL is an `Exporter` interface that backs up a MySQL database via the `mysqldump` command
type MySQL struct {
  // DB Host (e.g. 127.0.0.1)
  Host string
  // DB Port (e.g. 3306)
  Port string
  // DB Name
  DB string
  // DB User
  User string
  // DB Password
  Password string
}

var mysqlDump = func(x MySQL, path string) error {
  args := []string{
    fmt.Sprintf(`-r%v`, path),
    fmt.Sprintf(`-h%v`, x.Host),
    fmt.Sprintf(`-P%v`, x.Port),
    fmt.Sprintf(`-u%v`, x.User),
    fmt.Sprintf(`-p%v`, x.Password),
    x.DB,
  }

  _, err := exec.Command("mysqldump", args...).Output()
  return err
}

var mysqlTar = func(x MySQL, path string, destPath string) error {
  _, err := exec.Command("tar", "-cz", "-f"+destPath, path).Output()
  if (err != nil) { return err }

  err = os.Remove(path)
  return err
}

// Produces a `mysqldump` of the specified database, and creates a gzip compressed tarball archive.
func (x MySQL) Export() (*ExportResult) {
  result := &ExportResult{MIME: "application/x-tar"}

  dumpPath := fmt.Sprintf(`bu_%v_%v.sql`, x.DB, time.Now().Unix())

  result.Error = mysqlDump(x, dumpPath)
  if (result.Error != nil) { return result }

  tarPath := dumpPath+".tar.gz"
  result.Error = mysqlTar(x, dumpPath, tarPath)
  if (result.Error != nil) { return result }

  result.Path = tarPath
  return result
}

