package config

import (
	"reflect"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	t.Helper()

	conf, err := Get()
	if err != nil {
		t.Error(err)
	}

	expConfig := &Config{
		Timeout: 1 * time.Minute,
		TempDB:  "/tmp/dnscrypt",
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
			AllowedNames: "/etc/dnscrypt-proxy/allowed-names.txt",
			BlockedNames: "/etc/dnscrypt-proxy/blacklist-domains.txt",
		},
		SourceList: &Sources{
			Targets: []*RawTarget{
				&RawTarget{
					URL:    "https://raw.githubusercontent.com/PolishFiltersTeam/KADhosts/master/KADomains.txt",
					File:   "",
					Format: "domain",
					Type:   "blacklist",
					Notes:  "",
				},
				&RawTarget{
					URL:    "https://malc0de.com/bl/ZONES",
					File:   "",
					Format: "bind",
					Type:   "blacklist",
					Notes:  "",
				},
				&RawTarget{
					URL:    "https://urlhaus.abuse.ch/downloads/text",
					File:   "",
					Format: "url",
					Type:   "blacklist",
					Notes:  "",
				},
				&RawTarget{
					URL:    "https://joewein.net/dl/bl/dom-bl.txt",
					File:   "",
					Format: "host",
					Type:   "whitelist",
					Notes:  "public list of popular services",
				},
				&RawTarget{
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
