package target

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/maxkulish/dnscrypt-list/lib/config"
	"net/url"
	"sync"
)

// Format types of targets
type Format int8

const (
	// Domain example.com
	Domain Format = iota + 1
	// Host example.com
	Host
	// Bind zone "example.com"  {type master; file "/etc/hosts";};
	Bind
	// URL http://1.1.1.1:4/bin.sh
	URL
)

// Type of targets
type Type int

const (
	// WhiteList allowed targets
	WhiteList Type = iota + 1
	// BlackList targets to block
	BlackList
)

// Target keeps information about a target
type Target struct {
	TargetID   string
	RawURL     string
	TargetType Type
	URL        *url.URL
	Path       string
	TempFile   string
	Format     Format
	Notes      string
	Data       bytes.Buffer
	*sync.Mutex
}

// TypeString returns string representation of target Type
// t.TypeString(1) -> whitelist
// t.TypeString(2) -> blacklist
// t.TypeString(0) -> ""
// t.TypeString(5) -> ""
func (t *Target) TypeString() string {
	switch t.TargetType {
	case 1:
		return "whitelist"
	case 2:
		return "blacklist"
	default:
		return ""
	}
}

// NormalizeURL converts RawURL to URL
func (t *Target) NormalizeURL() {
	if normURL, err := url.Parse(t.RawURL); err == nil {
		t.URL = normURL
	}
}

// NewTargetFromRaw creates Target instance
// NewTargetFromRaw("blacklist", &rawTarget) -> *Target
func NewTargetFromRaw(tempDir string, rawTarget *config.RawTarget) *Target {

	// prepare local file name: tempDir + TargetType + Host
	// example: /tmp/dnscrypt/2_22_rescure.me

	t := &Target{
		TargetID:   uuid.New().String(),
		RawURL:     rawTarget.URL,
		Path:       rawTarget.File,
		Format:     getFormat(rawTarget.Format),
		TargetType: getType(rawTarget.Type),
	}
	t.NormalizeURL()
	t.TempFile = fmt.Sprintf("%s/%s", tempDir, t.TargetID)

	return t
}

// getFormat converts string to the Format type
// "domain" -> 1
// "host" -> 2
// "unknown" -> 0
func getFormat(format string) Format {
	switch format {
	case "domain":
		return Domain
	case "host":
		return Host
	case "bind":
		return Bind
	case "url":
		return URL
	default:
		return 0
	}
}

// getType
// "whitelist" -> 1
// "blacklist" -> 2
// "unknown" -> 0
func getType(rawType string) Type {
	switch {
	case rawType == "whitelist" || rawType == "white" || rawType == "allow":
		return WhiteList
	case rawType == "blacklist" || rawType == "black" || rawType == "block":
		return BlackList
	default:
		return 0
	}
}

// StringToBase64 produce base64 encoding for URL or Path
// string -> string
// StringToBase64("hello") -> "aGVsbG8="
func StringToBase64(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}
