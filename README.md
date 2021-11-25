# dnscrypt-list

CLI which generates blocked-names, blocked-ips, allowed-names for dnscrypt-proxy

### Types of sources

`url` - http or https links to the list of domains

`file` - path to the local file

### Formats

`domain` - new domain each line

`bind` - BIND zone file. Example: `zone "domain.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};`

`host` - similar to `/etc/hosts` Example: `0.0.0.0 domain.com`

### Temporary directory

Folder for Badger DB will be created after the first start

```shell
/tmp/dnscrypt/
├── 000001.vlog
├── DISCARD
├── KEYREGISTRY
└── MANIFEST
```

You can define another path in `config.yml`

## Build

### Local development

To build local version you need to install `goreleaser`

macOS

```shell
brew install goreleaser/tap/goreleaser
brew install goreleaser
```

Build binary for the current OS

```shell
make build-local
```

output

```shell
dist
├── config.yaml
└── dnscrypt-list_darwin_amd64
    └── dnscrypt-list

```

To run the binary with the `config.yml` in the current folder

```shell
./dist/dnscrypt-list_darwin_amd64/dnscrypt-list
```

or specify path to the config file

```shell
./dist/dnscrypt-list_darwin_amd64/dnscrypt-list --config=./dir/config.yml
```
