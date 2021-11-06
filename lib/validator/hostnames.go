package validator

import (
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

var (
	// ErrorHost can't check this host
	ErrorHost = errors.New("host checking error")
	// ErrorScheme domain wihtout scheme: http or https
	ErrorScheme = errors.New("scheme is obligatory")
)

// IsURL defines is it URL or local system path
// http://example.com -> true
// /etc/hosts -> false
func IsURL(urlOrPath string) (bool, error) {

	parse, err := url.Parse(urlOrPath)
	if err != nil {
		logger.Debug("url parsing error", zap.Error(err))
		return false, err
	}

	// check scheme
	if !parse.IsAbs() {
		return false, ErrorScheme
	}

	// check host
	if parse.Host != "" {
		return IsValidHost(parse.Host), nil
	}

	return false, ErrorHost
}

// IsValidHost checks if host is FQDN
// google.com -> true
// google -> false
func IsValidHost(host string) bool {

	if host == "" {
		return false
	}

	allString := hostnameRegexRFC952.FindString(host)

	if len(allString) > 0 {
		return true
	}
	return false

}

// IsHostReachable checks if the host is reachable by sending tcp request
// google.com -> true, nil
// asdggatil -> false, error
func IsHostReachable(host string) (bool, error) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		logger.Error("http.NewRequest error", zap.Error(err))
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http.Get url error", zap.Error(err))
		return false, err
	}
	defer resp.Body.Close()

	return true, nil
}
