package main

import (
	"testing"

	aptly_package_copier "github.com/bborbe/aptly_utils/package_copier"
	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()

	package_copier := aptly_package_copier.New(nil, nil, nil)

	err = do(writer, package_copier, "", "", "", "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
