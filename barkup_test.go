package barkup

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func testServer(code int, body string, contentType string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", contentType)
		fmt.Fprintln(w, body)
	}))
}

//////

func Test_ExportRestult_To_Move(t *testing.T) {
	file, _ := os.Create("to_mv_test")
	defer file.Close()

	e := ExportResult{"to_mv_test", "text/plain", nil}
	storeErr := e.To("test/", nil)
	expect(t, storeErr, (*Error)(nil))

	err := os.Remove("test/to_mv_test")
	expect(t, err, nil)
}

type StoreSuccessStory struct{}

func (x *StoreSuccessStory) Store(r *ExportResult, d string) *Error {
	return nil
}

type StoreFailureStory struct{}

func (x *StoreFailureStory) Store(r *ExportResult, d string) *Error {
	return &Error{
		err: errors.New("*****"),
	}
}

func Test_ExportRestult_To_Store(t *testing.T) {
	_, _ = os.Create("test/test.txt")
	e := &ExportResult{"test/test.txt", "text/plain", nil}
	err := e.To("test/", &StoreSuccessStory{})
	expect(t, err, (*Error)(nil))
}

func Test_ExportRestult_To_Store_Fail(t *testing.T) {
	_, _ = os.Create("test/test.txt")
	e := &ExportResult{"test/test.txt", "text/plain", nil}
	err := e.To("test/", &StoreFailureStory{})
	refute(t, err, (*Error)(nil))
}

/// Error

func Test_Error(t *testing.T) {
	e := Error{
		err:       fmt.Errorf("cheese"),
		CmdOutput: "CHEEEEES",
	}
	expect(t, e.Error(), "cheese")
}
