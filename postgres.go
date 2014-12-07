package barkup

import (
  "fmt"
  "time"
  "os/exec"
  "os"
)

// Postgres is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command
type Postgres struct {
  // DB Host (e.g. 127.0.0.1)
  Host string
  // DB Port (e.g. 5432)
  Port string
  // DB Name
  DB string
  // Connection Username
  Username string
  // Extra pg_dump options
  // e.g []string{"--inserts"}
  Options []string
}

// Produces a `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x Postgres) Export() (*ExportResult) {
  result := &ExportResult{MIME: "application/x-tar"}

  dumpPath := fmt.Sprintf(`bu_%v_%v.sql`, x.DB, time.Now().Unix())

  result.Error = pgDump(x, dumpPath)
  if (result.Error != nil) { return result }

  tarPath := dumpPath+".tar.gz"
  result.Error = pgTar(x, dumpPath, tarPath)
  if (result.Error != nil) { return result }

  result.Path = tarPath
  return result
}

func (x Postgres) dumpOptions() []string {
  options := x.Options

  if x.DB != "" {
    options = append(options, fmt.Sprintf(`-d%v`, x.DB))
  }

  if x.Host != "" {
    options = append(options, fmt.Sprintf(`-h%v`, x.Host))
  }

  if x.Port != "" {
    options = append(options, fmt.Sprintf(`-p%v`, x.Port))
  }

  if x.Username != "" {
    options = append(options, fmt.Sprintf(`-U%v`, x.Username))
  }

  return options
}

var pgDump = func(x Postgres, path string) error {
  options := append(x.dumpOptions(), fmt.Sprintf(`-f%v`, path))
  _, err := exec.Command("pg_dump", options...).Output()
  return err
}

var pgTar = func(x Postgres, path string, destPath string) error {
  _, err := exec.Command("tar", "-cz", "-f"+destPath, path).Output()
  if (err != nil) { return err }
  return os.Remove(path)
}