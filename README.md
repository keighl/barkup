# Barkup

[![Build Status](https://travis-ci.org/keighl/barkup.png?branch=master)](
https://travis-ci.org/keighl/barkup) [![Coverage Status](https://coveralls.io/repos/keighl/barkup/badge.svg?branch=master)](https://coveralls.io/r/keighl/barkup?branch=master)

[godoc.org/github.com/keighl/barkup](http://godoc.org/github.com/keighl/barkup)

Barkup is a library for backing things up. It provides tools for writing bare-bones backup programs in Go. The library is broken out into **exporters** and **storers**. Currently, those are:

**Exporters:** `MySQL` `Postgres` `RethinkDB`

**Storers:** `S3`

## Quick Example (mysql to s3)

Here's a go program that backups up a MySQL database (`Exporter`) to an S3 bucket (`Storer`) using barkup. The resulting executable is plopped on a server somewhere and scheduled to execute via CRON.

```go
package main

import "github.com/keighl/barkup"

func main() {

    // Configure a MySQL exporter
    mysql := &barkup.MySQL{
        Host:     "localhost",
        Port:     "3306",
        DB:       "production_db",
        User:     "root",
        Password: "cheese",
    }

    // Configure a S3 storer
    s3 := &barkup.S3{
        Region:       "us-east-1",
        Bucket:       "backups",
        AccessKey:    "*************",
        ClientSecret: "**********************",
    }

    // Export the database, and send it to the
    // bucket in the `db_backups` folder
    err := mysql.Export().To("db_backups/", s3)
    if err != nil {
        panic(err)
    }
}
```

```
$ go build
```

```
@hourly /path/to/backup-program
```

## Cron

Because barkup executes system commands like `tar`, `mysqldump`, etc, you may run into issues where crontab can't find 'em. The crontab shell is super stripped down, but you can shim the `PATH` environment variable to access stuff.

```
PATH=$PATH:/usr/local/bin:/usr/bin:/bin

@hourly /path/to/backup-program
```

This is especially the case for the RethinkDB exporter `rethink-dump` which executes system commands itself ("exec inception").

* * * * * 

## Exporters

Exporters provide a common interface for backing things up via the `Export()` method. It writes an export file to the local disk, and returns an `ExportResult` which can be passed on to a [storer](#storers), or to another location on the disk.

```go
// Exporter
mysql := &barkup.MySQL{...}

// Export Result
result := mysql.Export()
if (result.Error != nil) { panic(result.Error) }

// Send it to a directory path on a storer
err := result.To("backups/", storer)

// OR just move it somewhere on the local disk
err := result.To("~/backups/", nil)
```

---

### MySQL

The mysql exporter uses `mysqldump` to pump out a gzipped archive of your database. `mysqldump` must be installed on your system (it probably is if you're using mysql), and accessible to the user running the final program (again, it probabaly is).

**Usage**

```go
mysql := &barkup.MySQL{
  Host: "127.0.0.1",
  Port: "3306",
  DB: "XXXXX",
  User: "XXXXX",
  Password: "XXXXX",
  // Any extra mysqldump options
  Options: []string{"--skip-extended-insert"}
}

// Writes a file `./bu_DBNAME_TIMESTAMP.sql.tar.gz`
result := mysql.Export()

if (result.Error != nil) { panic(result.Error) }
```
---

### Postgres

The postgres exporter uses `pg_dump` to make a gzipped archive of your database. `pg_dump` must be installed on your system (it probably is if you're using postgres).

**Usage**

```go
postgres := &barkup.Postgres{
  Host: "127.0.0.1",
  Port: "5432",
  DB: "XXXXXXXX",

  // Not necessary if the program runs as an authorized pg user/role
  Username: "XXXXXXXX",

  // Any extra pg_dump options
  Options: []string{"--no-owner"},
}

// Writes a file `./bu_DBNAME_TIMESTAMP.sql.tar.gz`
result := postgres.Export()

if (result.Error != nil) { panic(result.Error) }
```

**Connection credentials**

You have two options for allowing barkup to connect to your DB.

Add a [`~/.pgpass`](http://www.postgresql.org/docs/9.3/static/libpq-pgpass.html) for account that will run the backup program.

Or, run the backup program from an authenticated user, like `postgres`:

```bash
$ sudo -i -u postgres
$ ./backup-program
```
---

### RethinkDB

The RethinkDB exporter uses `rethinkdb dump` to make a gzipped archive of your cluster. `rethinkdb-dump` must be installed on your system.

`$ sudo pip install rethinkdb`

**Usage**

```go
rethink := &barkup.RethinkDB{
  Name: "nightly",
  Connection: "0.0.0.0:28015",
  // You can specify specific databases and/or tables to dump (by default it dumps your whole cluster)
  Targets: []string{"site", "leads.contacts"},
}

// Writes a file `./bu_nightly_TIMESTAMP.tar.gz`
result := rethink.Export()
if (result.Error != nil) { panic(result.Error) }
```

## Storers

Storers take an `ExportResult` object and provide a common interface for moving a backup to someplace safe.

```go
// For chaining an ExportRestult
err := someExportResult.To("backups/", someStorer)

// OR
err := someStorer.Store(someExportResult, "backups/")
```

---

### S3

The S3 storer puts the exported file into a bucket at a specified directory. **Note,** you shouldn't use your global AWS credentials for this. Instead, [create bucket specific credentials via IAM.](http://blogs.aws.amazon.com/security/post/Tx3VRSWZ6B3SHAV/Writing-IAM-Policies-How-to-grant-access-to-an-Amazon-S3-bucket)

**Usage**

```go
s3 := &barkup.S3{
  Region: "us-east-1",
  Bucket: "backups",
  AccessKey: "XXXXXXXXXXXXX",
  ClientSecret: "XXXXXXXXXXXXXXXXXXXXX",
}

err := someExportResult.To("data/", s3)
```

**Region IDs**

* us-east-1
* us-west-1
* us-west-2
* eu-west-1
* ap-southeast-1
* ap-southeast-2
* ap-northeast-1
* sa-east-1

