package barkup

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	// TarCmd is the path to the `tar` executable
	TarCmd = "tar"
	// MysqlDumpCmd is the path to the `mysqldump` executable
	MysqlDumpCmd = "mysqldump"
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
	// Extra mysqldump options
	// e.g []string{"--extended-insert"}
	Options []string
}

// Export produces a `mysqldump` of the specified database, and creates a gzip compressed tarball archive.
func (x MySQL) Export() *ExportResult {
	result := &ExportResult{MIME: "application/x-tar"}

	dumpPath := fmt.Sprintf(`bu_%v`, time.Now().Unix())
	var (
		dbs []string
		err error
	)
	if x.DB != "" {
		dbs = []string{x.DB}
	} else {
		dbs, err = x.getDBNames()
		if err != nil {
			result.Error = makeErr(err, "MySQL connection error")
			return result
		}
	}

	os.Mkdir(dumpPath, 0770)
	for _, db := range dbs {
		dumpFile := fmt.Sprintf(`bu_%v_%v.sql`, db, time.Now().Unix())
		options := append(x.dumpOptions(db), fmt.Sprintf(`-r%v/%v`, dumpPath, dumpFile))
		out, err := exec.Command(MysqlDumpCmd, options...).Output()
		if err != nil {
			result.Error = makeErr(err, string(out))
			return result
		}
	}

	result.Path = dumpPath + ".tar.gz"
	out, err := exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		return result
	}
	os.RemoveAll(dumpPath)

	return result
}

func (x MySQL) getDBNames() ([]string, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/", x.User, x.Password, x.Host, x.Port))
	if err != nil {
		return []string{}, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT SCHEMA_NAME FROM `information_schema`.`SCHEMATA` WHERE SCHEMA_NAME NOT REGEXP 'information_schema|performance_schema|sys'")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var dbs []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		dbs = append(dbs, name)
	}
	return dbs, nil
}

func (x MySQL) dumpOptions(db string) []string {
	options := x.Options
	options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	options = append(options, fmt.Sprintf(`-P%v`, x.Port))
	options = append(options, fmt.Sprintf(`-u%v`, x.User))
	options = append(options, fmt.Sprintf(`-p%v`, x.Password))
	options = append(options, db)
	return options
}
