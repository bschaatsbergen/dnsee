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
A       google.com.     146     142.251.39.110
AAAA    google.com.     300     2a00:1450:400e:803::200e
MX      google.com.     193     smtp.google.com.        10
NS      google.com.     103     ns1.google.com.
NS      google.com.     103     ns4.google.com.
NS      google.com.     103     ns3.google.com.
NS      google.com.     103     ns2.google.com.
```

### Fetch all records for a specific type

To get all records for a domain name of a specific type:

```
$ dnsee google.com -q A
A       google.com.     146     142.251.39.110
```

Multiple types can be passed at once:
```
$ dnsee google.com -q 'A SOA'
A	google.com.	105	142.251.36.46
SOA	google.com.	20	ns1.google.com.	dns-admin.google.com.
```

### Fetch all records using a different DNS server

To get all records for a domain name using a different DNS server:

```
$ dnsee google.com --dns-server-ip 1.1.1.1
A       google.com.     146     142.251.39.110
AAAA    google.com.     300     2a00:1450:400e:803::200e
MX      google.com.     193     smtp.google.com.        10
NS      google.com.     103     ns1.google.com.
NS      google.com.     103     ns4.google.com.
NS      google.com.     103     ns3.google.com.
NS      google.com.     103     ns2.google.com.
```

## Contributing

Contributions are highly appreciated and always welcome.
Have a look through existing [Issues](https://github.com/bschaatsbergen/dnsee/issues) and [Pull Requests](https://github.com/bschaatsbergen/dnsee/pulls) that you could help with.
