package barkup

import (
	"fmt"
	"os/exec"
	"time"
)

var (
	// RethinkCmd is the path to the `rethinkdb` executable
	RethinkCmd = "rethinkdb"
)

// RethinkDB is an `Exporter` interface that backs up a RethinkDB cluster via the `rethinkdb dump` command
type RethinkDB struct {
	// Describe the dataset you're backing up. Used in filenaming
	// e.g. website
	Name string
	// HOST:PORT Connection (e.g. 127.0.0.1:28015)
	Connection string
	// Limit the dump to specific databases or tables in the cluster
	// e.g. `[]string{"main_website", "internal_website.users"}`
	Targets []string
	// Extra "rethinkdb dump" options
	// See `rethinkdb dump -h`
	// e.g  `[]string{"--authAUTHKEY"}`
	Options []string
}

// Export produces a gzip compressed tarball archive of the rethink cluster (or targetted DBs/tables)
func (x RethinkDB) Export() *ExportResult {
	result := &ExportResult{MIME: "application/x-tar"}
	result.Path = fmt.Sprintf(`bu_%v_%v.tar.gz`, x.Name, time.Now().Unix())
	options := append(x.dumpOptions(), fmt.Sprintf(`-f%v`, result.Path))
	out, err := exec.Command(RethinkCmd, options...).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
	}
	return result
}

func (x RethinkDB) dumpOptions() []string {
	options := append([]string{"dump"}, x.Options...)
	options = append(options, fmt.Sprintf(`-c%v`, x.Connection))
	for _, target := range x.Targets {
		options = append(options, fmt.Sprintf(`-e%v`, target))
	}
	return options
}
