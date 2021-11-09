package files

import (
	"fmt"
	"os"
	"testing"
)

const (
	testFile = "./test.txt"
)

func TestIsIsPathExist(t *testing.T) {
	t.Helper()

	err := CreateFileIfNotExist(testFile)
	if err != nil {
		t.Error("can't create test file", err)
	}
	defer func() {
		err = os.RemoveAll(testFile)
	}()

	pathValid := func(path string, expRes bool) {
		res := IsPathExist(path)

		if res != expRes {
			t.Fatalf("unexpected result for IsValidPath(%s); got: %t; want: %t", path, res, expRes)
		}
	}

	pathValid(testFile, true)
	pathValid("/opt/dnscrypt-proxy/whitelist-private.txt", true)
	pathValid("/abc/unknown.txt", false)
}

func ExampleIsPathExist() {
	fmt.Println(IsPathExist("/opt/dnscrypt-proxy/whitelist-private.txt"))
	fmt.Println(IsPathExist("/abc/unknown.txt"))
	// Output:
	// true
	// false
}

func BenchmarkIsPathExist(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		IsPathExist("/opt/dnscrypt-proxy/whitelist-private.txt")
	}
}
