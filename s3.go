package barkup

import (
	"bufio"
	"github.com/go-amz/amz/aws"
	"github.com/go-amz/amz/s3"
	"os"
)

// S3 is a `Storer` interface that puts an ExportResult to the specified S3 bucket. Don't use your main AWS keys for this!! Create read-only keys using IAM
type S3 struct {
	// Available regions:
	// * us-east-1
	// * us-west-1
	// * us-west-2
	// * eu-west-1
	// * ap-southeast-1
	// * ap-southeast-2
	// * ap-northeast-1
	// * sa-east-1
	Region string
	// Name of the bucjet
	Bucket string
	// AWS S3 access key
	AccessKey string
	// AWS S3 secret
	ClientSecret string
}

// Store puts an `ExportResult` struct to an S3 bucket within the specified directory
func (x *S3) Store(result *ExportResult, directory string) *Error {

	if result.Error != nil {
		return result.Error
	}

	file, err := os.Open(result.Path)
	if err != nil {
		return makeErr(err, "")
	}
	defer file.Close()

	buffy := bufio.NewReader(file)
	stat, err := file.Stat()
	if err != nil {
		return makeErr(err, "")
	}

	size := stat.Size()

	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}

	s := s3.New(auth, aws.Regions[x.Region])
	bucket := s.Bucket(x.Bucket)

	err = bucket.PutReader(directory+result.Filename(), buffy, size, result.MIME, s3.BucketOwnerFull)
	return makeErr(err, "")
}
