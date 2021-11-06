package validator

import "regexp"

const (
	hostnameRegexStringRFC952 = `^[a-zA-Z][a-zA-Z0-9\-\.]+[a-z-Az0-9]$` // https://tools.ietf.org/html/rfc952
)

var (
	hostnameRegexRFC952 = regexp.MustCompile(hostnameRegexStringRFC952)
)
