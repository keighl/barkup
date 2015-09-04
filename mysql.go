package barkup

import (
  "fmt"
  "time"
  "os/exec"
  "os"
)

// TarCmd is the path to the `tar` executable
var TarCmd string = "tar"

// MysqlDumpCmd is the path to the `mysqldump` executable
var MysqlDumpCmd string = "mysqldump"

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
  // Extra mysqldump options
  // e.g []string{"--extended-insert"}
  Options []string
}

// Produces a `mysqldump` of the specified database, and creates a gzip compressed tarball archive.
func (x MySQL) Export() (*ExportResult) {
  result := &ExportResult{MIME: "application/x-tar"}

  dumpPath := fmt.Sprintf(`bu_%v_%v.sql`, x.DB, time.Now().Unix())

  options := append(x.dumpOptions(), fmt.Sprintf(`-r%v`, dumpPath))
  _, err := exec.Command(MysqlDumpCmd, options...).Output()
  result.Error = err
  if (result.Error != nil) { return result }

  result.Path = dumpPath+".tar.gz"
  _, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()
  result.Error = err
  if (err != nil) { return result }
  os.Remove(dumpPath)

  return result
}

func (x MySQL) dumpOptions() []string {
  options := x.Options
  options = append(options, fmt.Sprintf(`-h%v`, x.Host))
  options = append(options, fmt.Sprintf(`-P%v`, x.Port))
  options = append(options, fmt.Sprintf(`-u%v`, x.User))
  options = append(options, fmt.Sprintf(`-p%v`, x.Password))
  options = append(options, x.DB)
  return options
}