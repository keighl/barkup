package barkup

import (
	"bufio"
	"os"

	"github.com/ncw/swift"
)

// ObjectStorage is a `Storer` interface that puts an ExportResult to the specified container/bucket. This can be used to backup to rackspace/swift/openstack and softlayers object storage solutions
type ObjectStorage struct {
	APIKey      string
	Username    string
	Container   string
	AuthURL     string
	Hash        string
	ContentType string
	CheckHash   bool
	Headers     swift.Headers
}

// Store stores the result in the given directory path for the container specified in the caller
func (o *ObjectStorage) Store(result *ExportResult, directory string) *Error {
	if result.Error != nil {
		return result.Error
	}

	file, err := os.Open(result.Path)
	if err != nil {
		return makeErr(err, "")
	}
	defer file.Close()

	buffy := bufio.NewReader(file)

	// Create a v1 auth connection
	c := swift.Connection{
		UserName: o.Username,
		ApiKey:   o.APIKey,
		AuthUrl:  o.AuthURL,
	}

	// Authenticate
	err = c.Authenticate()
	if err != nil {
		makeErr(err, "")
	}

	_, err = c.ObjectPut(o.Container, directory, buffy, o.CheckHash, o.Hash, o.ContentType, o.Headers)
	return makeErr(err, "")
}
