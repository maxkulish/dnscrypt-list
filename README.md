# dnscrypt-list
CLI which generates blocked-names, blocked-ips, allowed-names for dnscrypt-proxy


### Types of sources
`url` - http or https links to the list of domains

`file` - path to the local file

### Formats
`domain` - new domain each line

`bind` - BIND zone file.  Example: `zone "domain.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};` 

`host` - similar to `/etc/hosts` Example: `0.0.0.0 domain.com`