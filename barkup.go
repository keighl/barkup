package barkup

import (
  "os/exec"
  "path/filepath"
)

//////////////

// The `Export()` method is expected to export "something" to a file and return a complete `ExportResult` struct (`Path`, `MIME`, `Error`). If any error occurs during it's work, it should set the error to the result's `Error` attribute
type Exporter interface {
  Export() (*ExportResult, error)
}

//////////////

// The `Store()` method takes an `ExportResult` and move it somewhere! To a cloud storage service, for instance...
type Storer interface {
  Store(result *ExportResult, directory string) (error)
}

//////////////

type ExportResult struct {
  // Path to exported file
  Path string
  // MIME type of the exported file (e.g. application/x-tar)
  MIME string
  // Any error that occured during `Export()`
  Error error
}

// Hands off an ExportResult to a `Storer` interface and invokes its Store() method. The directory argument is passed along too. If `store` is `nil`, the the method will simply move the export result to the specified directory (via the `mv` command)
func (x *ExportResult) To(directory string, store Storer) (error) {
  if (store == nil) {
    _, err := exec.Command("mv", x.Path, directory + x.Filename()).Output()
    return err
  }

  err := store.Store(x, directory)
  return err
}

// Returns the just filename component of the `Path` attribute
func (x ExportResult) Filename() (string) {
  _, filename := filepath.Split(x.Path)
  return filename
}


