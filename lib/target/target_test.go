package target

import (
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"testing"
)

func init() {
	logger.SetLogger()
}

func TestNewTarget(t *testing.T) {

}

func TestStringToBase64(t *testing.T) {

	t.Helper()

	strToBase64 := func(str, expRes string) {
		res := StringToBase64(str)

		if res != expRes {
			t.Fatalf("unexpected result for StringToBase64(%v); got: %s; want: %s", str, res, expRes)
		}

	}

	strToBase64("", "")
	strToBase64("hello", "aGVsbG8=")
	strToBase64("/tmp/dnscrypt/blacklist-domains.txt", "L3RtcC9kbnNjcnlwdC9ibGFja2xpc3QtZG9tYWlucy50eHQ=")
	strToBase64("github.com/maxkulish/dnscrypt-list", "Z2l0aHViLmNvbS9tYXhrdWxpc2gvZG5zY3J5cHQtbGlzdA==")
	strToBase64("https://github.com/maxkulish/dnscrypt-list", "aHR0cHM6Ly9naXRodWIuY29tL21heGt1bGlzaC9kbnNjcnlwdC1saXN0")
}
