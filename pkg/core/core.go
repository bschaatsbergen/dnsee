package core

import (
	"fmt"
	"net"
	"runtime"
	"time"

	model "github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/fatih/color"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const windows = "windows"
const resolverPath = "/etc/resolv.conf"

// GetQueryTypes returns a slice of all supported DNS query types.
func GetQueryTypes() []model.QueryType {
	return []model.QueryType{
		{Type: dns.TypeA, Name: "A"},
		{Type: dns.TypeAAAA, Name: "AAAA"},
		{Type: dns.TypeCNAME, Name: "CNAME"},
		{Type: dns.TypeMX, Name: "MX"},
		{Type: dns.TypeTXT, Name: "TXT"},
		{Type: dns.TypeNS, Name: "NS"},
		{Type: dns.TypeSOA, Name: "SOA"},
		{Type: dns.TypePTR, Name: "PTR"},
	}
}

// FilterQueryTypes filters the queryTypes slice to only include the query type specified by the user.
func FilterQueryTypes(queryTypes []model.QueryType, userSpecifiedQueryType string) []model.QueryType {
	var filteredQueryTypes []model.QueryType
	for _, queryType := range queryTypes {
		if queryType.Name == userSpecifiedQueryType {
			filteredQueryTypes = append(filteredQueryTypes, queryType)
			break
		}
	}
	return filteredQueryTypes
}

// PrepareDNSQuery prepares a DNS query for a given domain name and query type.
func PrepareDNSQuery(domainName string, queryType uint16) dns.Msg {
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domainName), queryType)
	return msg
}

// SendDNSQuery sends a DNS query to a given DNS server.
func SendDNSQuery(client *dns.Client, msg dns.Msg, dnsServerIP string) (*dns.Msg, time.Duration, error) {
	if dnsServerIP == "" {
		goOS := runtime.GOOS
		if goOS == windows {
			logrus.Fatal("error: Unable to retrieve DNS configuration on Windows. \nPlease specify a DNS server IP explicitly with the `--dns-server-ip` flag.")
		}
		conf, err := dns.ClientConfigFromFile(resolverPath)
		if err != nil {
			logrus.Errorf("error: %s. Unable to retrieve DNS configuration.", err)
			logrus.Fatal("Please specify a DNS server IP explicitly with the `--dns-server-ip` flag.")
		}
		dnsServerIP = conf.Servers[0]

	}
	// If the server IP is IPv6, wrap it in square brackets to remove ambiguity from port number and address.
	if net.ParseIP(dnsServerIP).To4() == nil {
		dnsServerIP = "[" + dnsServerIP + "]"
	}

	logrus.Debugf("Sending DNS query to %s", dnsServerIP)
	response, timeDuration, err := client.Exchange(&msg, dnsServerIP+":53")

	if err != nil {
		logrus.Debug("Failed to receive DNS response.")
		logrus.Debugf("Query sent to DNS server:\n%s", msg.String())
		logrus.Fatal(err)
	}
	logrus.Debugf("Received DNS response from %s, Round-trip time: %s", dnsServerIP, timeDuration)
	logrus.Debugf("Response body:\n%s", response.String())
	return response, timeDuration, nil
}

// DisplayRecords displays the DNS records returned by the DNS server.
func DisplayRecords(domainName string, queryType struct {
	Type uint16
	Name string
}, answers []dns.RR) {
	for _, ans := range answers {
		switch queryType.Type {
		case dns.TypeA:
			if aRecord, ok := ans.(*dns.A); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), aRecord.Hdr.Ttl, aRecord.A)
			}
		case dns.TypeAAAA:
			if aaaaRecord, ok := ans.(*dns.AAAA); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), aaaaRecord.Hdr.Ttl, aaaaRecord.AAAA)
			}
		case dns.TypeCNAME:
			if cnameRecord, ok := ans.(*dns.CNAME); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), cnameRecord.Hdr.Ttl, cnameRecord.Target)
			}
		case dns.TypeMX:
			if mxRecord, ok := ans.(*dns.MX); ok {
				fmt.Printf("%s\t%s.\t%d\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), mxRecord.Hdr.Ttl, mxRecord.Preference, mxRecord.Mx)
			}
		case dns.TypeTXT:
			if txtRecord, ok := ans.(*dns.TXT); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), txtRecord.Hdr.Ttl, txtRecord.Txt[0])
			}
		case dns.TypeNS:
			if nsRecord, ok := ans.(*dns.NS); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), nsRecord.Hdr.Ttl, nsRecord.Ns)
			}
		case dns.TypeSOA:
			if soaRecord, ok := ans.(*dns.SOA); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), soaRecord.Hdr.Ttl, soaRecord.Ns, soaRecord.Mbox)
			}
		case dns.TypePTR:
			if ptrRecord, ok := ans.(*dns.PTR); ok {
				fmt.Printf("%s\t%s.\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), ptrRecord.Hdr.Ttl, ptrRecord.Ptr)
			}
		}
	}
}
