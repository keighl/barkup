package barkup

import (
	"os"
	"testing"

	"github.com/ncw/swift/swifttest"
)

func Test_Object_Storeage_Store_Success(t *testing.T) {
	stest, err := swifttest.NewSwiftServer("0.0.0.0")

	s := &ObjectStorage{
		Username:  swifttest.TEST_ACCOUNT,
		APIKey:    swifttest.TEST_ACCOUNT,
		Container: "test",
		AuthURL:   stest.AuthURL,
	}

	_, err = os.Create("test.txt")
	err = s.Store(&ExportResult{"test.txt", "text/plain", nil}, "")
	expect(t, err, (*Error)(nil))
}

func Test_Object_Storage_Store_Fail(t *testing.T) {
	stest, err := swifttest.NewSwiftServer("0.0.0.0")
	stest.Close()

	s := &ObjectStorage{
		Username:  swifttest.TEST_ACCOUNT,
		APIKey:    swifttest.TEST_ACCOUNT,
		Container: "test",
		AuthURL:   stest.AuthURL,
	}

	_, _ = os.Create("test.txt")
	err = s.Store(&ExportResult{"test.txt", "text/plain", nil}, "")
	refute(t, err, (*Error)(nil))
}

func Test_Object_Storage_Store_ExportError(t *testing.T) {
	stest, err := swifttest.NewSwiftServer("0.0.0.0")

	s := &ObjectStorage{
		Username:  swifttest.TEST_ACCOUNT,
		APIKey:    swifttest.TEST_ACCOUNT,
		Container: "test",
		AuthURL:   stest.AuthURL,
	}

	_, _ = os.Create("test/test.txt")
	err = s.Store(&ExportResult{"", "text/plain", &Error{}}, "test/")
	refute(t, err, nil)
}
