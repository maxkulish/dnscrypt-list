package validator

import (
	"fmt"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"testing"
)

func init() {
	logger.SetLogger()
}

func TestIsURL(t *testing.T) {
	t.Helper()

	isURL := func(URL string, expRes bool) {

		res, _ := IsURL(URL)

		if res != expRes {
			t.Fatalf("unexpected result for IsURL(%s); got: %t; want: %t", URL, res, expRes)
		}

	}

	isURL("http://google.com", true)
	isURL("https://google.com", true)
	isURL("http://google.local", true)
	isURL("google.com", false)
	isURL("/etc/hosts", false)
	isURL("localhost", false)
}

func ExampleIsURL() {
	fmt.Println(IsURL("http://google.com"))
	fmt.Println(IsURL("/etc/hosts"))
	// Output:
	// true <nil>
	// false scheme is obligatory
}

func TestIsValidHost(t *testing.T) {
	t.Helper()

	isURL := func(host string, expRes bool) {

		res := IsValidHost(host)

		if res != expRes {
			t.Fatalf("unexpected result for IsValidHost(%s); got: %t; want: %t", host, res, expRes)
		}

	}

	isURL("http://google.com", false)
	isURL("https://google.com", false)
	isURL("https://google.com:8080", false)
	isURL("http://google.local", false)
	isURL("google.com", true)
	isURL(".google.com", false)
	isURL("typical-hostname33.whatever.co.uk", true)
	isURL("typical_hostname33.whatever.co.uk.local", false)
	isURL("llanfairpwllgwyngyllgogerychwyrndro-bwyll-llantysiliogogogoch.info", true)
	isURL("/etc/hosts", false)
	isURL("localhost", true)
	isURL("ab", false)
}
