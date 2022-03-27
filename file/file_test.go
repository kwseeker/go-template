package file

import "testing"

func TestCreateFileWithParentDir(t *testing.T) {
	CreateFileWithParentDir("/tmp/go/test/test.txt")
}
