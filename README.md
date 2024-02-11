# dnsee

[![Release](https://github.com/bschaatsbergen/dnsee/actions/workflows/goreleaser.yaml/badge.svg)](https://github.com/bschaatsbergen/dnsee/actions/workflows/goreleaser.yaml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bschaatsbergen/dnsee) ![GitHub commits since latest release (by SemVer)](https://img.shields.io/github/commits-since/bschaatsbergen/dnsee/latest) [![Go Reference](https://pkg.go.dev/badge/github.com/bschaatsbergen/dnsee.svg)](https://pkg.go.dev/github.com/bschaatsbergen/dnsee) ![GitHub all releases](https://img.shields.io/github/downloads/bschaatsbergen/dnsee/total) 

See DNS configurations quickly

## Brew

To install dnsee using brew, simply do the below.

```sh
brew tap bschaatsbergen/dnsee
brew install dnsee
```

## Binaries

You can download the [latest binary](https://github.com/bschaatsbergen/dnsee/releases/latest) for Linux, MacOS, and Windows.

## Examples

Using `dnsee` is very simple.

### Fetch all records

To get all records for a domain name:

```
$ dnsee google.com
A       gooogle.com.    04m42s          142.251.36.4
AAAA    gooogle.com.    04m42s          2a00:1450:400e:800::2004
MX      gooogle.com.    04m42s          0       .
TXT     gooogle.com.    04m42s          v=spf1 -all
NS      gooogle.com.    01h48m35s       ns2.google.com.
NS      gooogle.com.    01h48m35s       ns3.google.com.
NS      gooogle.com.    01h48m35s       ns1.google.com.
NS      gooogle.com.    01h48m35s       ns4.google.com.
SOA     gooogle.com.    42s             ns1.google.com. dns-admin.google.com.
```

### Fetch all records for a specific type

To get all records for a domain name of a specific type:

```
$ dnsee google.com -q A
A       google.com.     03m15s          216.58.214.14
```

### Fetch all records using a different DNS server

To get all records for a domain name using a different DNS server:

```
$ dnsee google.com --dns-server-ip 1.1.1.1
A       google.com.     01m02s          142.250.179.174
AAAA    google.com.     47s             2a00:1450:400e:811::200e
MX      google.com.     34s             10      smtp.google.com.
NS      google.com.     90h42m12s       ns1.google.com.
NS      google.com.     90h42m12s       ns4.google.com.
NS      google.com.     90h42m12s       ns3.google.com.
NS      google.com.     90h42m12s       ns2.google.com.
SOA     google.com.     44s             ns1.google.com. dns-admin.google.com.
```

## Contributing

Contributions are highly appreciated and always welcome.
Have a look through existing [Issues](https://github.com/bschaatsbergen/dnsee/issues) and [Pull Requests](https://github.com/bschaatsbergen/dnsee/pulls) that you could help with.
