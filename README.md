# Barkup

[![Build Status](https://travis-ci.org/keighl/barkup.png?branch=master)](
https://travis-ci.org/keighl/barkup)

Barkup is a library for backing things up. It provides tools for writing bare-bones backup programs in Go. The library is broken out into **exporters** and **storers**. Currently, those are only:

**Exporters:** `MySQL`

**Storers:** `S3`

### E.g. MySQL to S3

    package main

    import (
      "github.com/keighl/barkup"
    )

    func main() {

      // Configure a MySQL exporter
      mysql := &barkup.MySQL{
        Host: "localhost",
        Port: "3306",
        DB: "production_db",
        User: "root",
        Password: "cheese",
      }

      // Configure a S3 storer
      s3 := &barkup.S3{
        Region: "us-east-1",
        Bucket: "backups",
        AccessKey: "*************",
        ClientSecret: "**********************",
      }

      // Export the database, and send it to the
      // bucket in the `db_backups` folder
      err := mysql.Export().To("db_backups/", s3)
      if (err != nil) { panic(err) }
    }

