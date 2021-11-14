package download

import (
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"testing"
)

func init() {
	logger.SetLogger()
}

func TestNormalizedDomain(t *testing.T) {
	t.Helper()

	normDomain := func(domain, expRes string, expError error) {

		normDomain, err := NormalizedDomain(domain)
		if err != expError {
			t.Fatalf("unexpected error for NormalizedDomain(%s); got: %v; want: %v", domain, err, expError)
		}

		if normDomain != expRes {
			t.Fatalf("unexpected result for NormalizedDomain(%s); got: %s; want: %s", domain, normDomain, expRes)
		}
	}

	normDomain("    example.com  ", "example.com", nil)
	normDomain("   #  example.com  ", "", ErrCommentedLine)
	normDomain(".example.com  ", "", ErrInvalidDomain)
	normDomain("http://example.com  ", "", ErrInvalidDomain)
	normDomain("# example.com  ", "", ErrCommentedLine)
}
