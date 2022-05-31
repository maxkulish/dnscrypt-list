package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetConfig(t *testing.T) {
	t.Helper()

	conf, err := Get()
	if err != nil {
		t.Error(err)
	}

	expConfig := &Config{
		Timeout:     1 * time.Minute,
		WhiteListDB: "/tmp/dnscrypt/whitelist.db",
		BlackListDB: "/tmp/dnscrypt/blacklist.db",
		Update: Update{
			Sources:   1 * time.Hour,
			Blacklist: 5 * time.Minute,
			Whitelist: 6 * time.Hour,
		},
		SourcesLink: SourcesLink{
			FilePath: "sources_test.yml",
			URL:      "",
		},
		Output: Output{
			Whitelist: "/etc/dnscrypt-proxy/allowed-names.txt",
			Blacklist: "/etc/dnscrypt-proxy/blacklist-domains.txt",
		},
		SourceList: &Sources{
			Targets: []*RawTarget{
				{
					URL:    "https://raw.githubusercontent.com/PolishFiltersTeam/KADhosts/master/KADomains.txt",
					File:   "",
					Format: "domain",
					Type:   "blacklist",
					Notes:  "",
				},
				{
					URL:    "https://malc0de.com/bl/ZONES",
					File:   "",
					Format: "bind",
					Type:   "blacklist",
					Notes:  "",
				},
				{
					URL:    "https://urlhaus.abuse.ch/downloads/text",
					File:   "",
					Format: "url",
					Type:   "blacklist",
					Notes:  "",
				},
				{
					URL:    "https://joewein.net/dl/bl/dom-bl.txt",
					File:   "",
					Format: "host",
					Type:   "whitelist",
					Notes:  "public list of popular services",
				},
				{
					URL:    "",
					File:   "/opt/dnscrypt-proxy/whitelist-private.txt",
					Format: "domain",
					Type:   "whitelist",
					Notes:  "my private domains to allow even if some of them are on the blacklist",
				},
			},
		},
	}

	if conf.Timeout != expConfig.Timeout {
		t.Fatalf("unexpected result for config.Get();\n got: %+v;\n want: %+v\n", conf.Timeout, expConfig.Timeout)
	}

	for i := range conf.SourceList.Targets {
		if !reflect.DeepEqual(conf.SourceList.Targets[i], expConfig.SourceList.Targets[i]) {
			t.Fatalf("unexpected result for config.Get();\n got: %+v;\n want: %+v\n", conf.SourceList, expConfig.SourceList)
		}
	}

	if conf.Output != expConfig.Output {
		t.Fatalf("unexpected result for config.Get();\n got: %+v;\n want: %+v\n", conf.Output, expConfig.Output)
	}

	if conf.SourcesLink != expConfig.SourcesLink {
		t.Fatalf("unexpected result for config.Get();\n got: %+v;\n want: %+v\n", conf.SourcesLink, expConfig.SourcesLink)
	}

	if conf.Update != expConfig.Update {
		t.Fatalf("unexpected result for config.Get();\n got: %+v;\n want: %+v\n", conf.Update, expConfig.Update)
	}
}

func TestConfigFromYAML(t *testing.T) {
	t.Helper()
	t.Parallel()

	want := &Config{
		Timeout:     1 * time.Minute,
		TempDir:     "/tmp/dnscrypt",
		BlackListDB: "/tmp/dnscrypt/blacklist.db",
		WhiteListDB: "/tmp/dnscrypt/whitelist.db",
		Update: Update{
			Sources:   1 * time.Hour,
			Blacklist: 5 * time.Minute,
			Whitelist: 6 * time.Hour,
		},
		SourcesLink: SourcesLink{
			FilePath: "sources.yml",
			URL:      "",
		},
		Output: Output{
			Whitelist: "/tmp/dnscrypt/allowed-names.txt",
			Blacklist: "/tmp/dnscrypt/blacklist-domains.txt",
		},
	}
	got, err := configFromYAML("../../testdata/config.yml")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

func TestSourceFromYAML(t *testing.T) {
	t.Parallel()

	want := &Sources{
		Targets: []*RawTarget{
			{
				URL:    "",
				File:   "/opt/dnscrypt-proxy/whitelist-private.txt",
				Format: "domain",
				Type:   "whitelist",
				Notes:  "my private domains to allow even if some of them are on the blacklist",
			},
			{
				URL:    "https://raw.githubusercontent.com/PolishFiltersTeam/KADhosts/master/KADomains.txt",
				File:   "",
				Format: "domain",
				Type:   "blacklist",
				Notes:  "",
			},
			{
				URL:    "https://raw.githubusercontent.com/notracking/hosts-blocklists/master/dnscrypt-proxy/dnscrypt-proxy.blacklist.txt",
				File:   "",
				Format: "url",
				Type:   "blacklist",
				Notes:  "notracking",
			},
		},
	}

	got, err := sourceFromYAML("../../testdata/sources.yml")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}
