package barkup

import (
	"os"
	"os/exec"
	"path/filepath"
)

//////////////

// Exporter is expected to export "something" to a file and return a complete `ExportResult` struct (`Path`, `MIME`, `Error`). If any error occurs during it's work, it should set the error to the result's `Error` attribute
type Exporter interface {
	Export() (*ExportResult, *Error)
}

// Error can ship a cmd output as well as the start interface. Useful for understanding why a system command (exec.Command) failed
type Error struct {
	err       error
	CmdOutput string
}

func (e Error) Error() string {
	return e.err.Error()
}

func makeErr(err error, out string) *Error {
	if err != nil {
		return &Error{
			err:       err,
			CmdOutput: out,
		}
	}
	return nil
}

//////////////

// Storer takes an `ExportResult` and move it somewhere! To a cloud storage service, for instance...
type Storer interface {
	Store(result *ExportResult, directory string) *Error
}

//////////////

// ExportResult is the result of an export operation... duh
type ExportResult struct {
	// Path to exported file
	Path string
	// MIME type of the exported file (e.g. application/x-tar)
	MIME string
	// Any error that occured during `Export()`
	Error *Error
}

// To hands off an ExportResult to a `Storer` interface and invokes its Store() method. The directory argument is passed along too. If `store` is `nil`, the the method will simply move the export result to the specified directory (via the `mv` command)
func (x *ExportResult) To(directory string, store Storer) *Error {
	if store == nil {
		out, err := exec.Command("mv", x.Path, directory+x.Filename()).Output()
		return makeErr(err, string(out))
	}

	storeErr := store.Store(x, directory)
	if storeErr != nil {
		return storeErr
	}

	err := os.Remove(x.Path)
	return makeErr(err, "")
}

// Filename returns the just filename component of the `Path` attribute
func (x ExportResult) Filename() string {
	_, filename := filepath.Split(x.Path)
	return filename
}
