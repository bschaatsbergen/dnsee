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
func SendDNSQuery(client *dns.Client, msg dns.Msg, dnsServerIP, dnsServerPort string) (*dns.Msg, time.Duration, error) {
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

	addr := net.JoinHostPort(dnsServerIP, dnsServerPort)

	logrus.Debugf("Sending DNS query to %s", addr)
	response, timeDuration, err := client.Exchange(&msg, addr)

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
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(aRecord.Hdr.Ttl), aRecord.A)
			}
		case dns.TypeAAAA:
			if aaaaRecord, ok := ans.(*dns.AAAA); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(aaaaRecord.Hdr.Ttl), aaaaRecord.AAAA)
			}
		case dns.TypeCNAME:
			if cnameRecord, ok := ans.(*dns.CNAME); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(cnameRecord.Hdr.Ttl), cnameRecord.Target)
			}
		case dns.TypeMX:
			if mxRecord, ok := ans.(*dns.MX); ok {
				fmt.Printf("%s\t%s.\t%s\t%d\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(mxRecord.Hdr.Ttl), mxRecord.Preference, mxRecord.Mx)
			}
		case dns.TypeTXT:
			if txtRecord, ok := ans.(*dns.TXT); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(txtRecord.Hdr.Ttl), txtRecord.Txt[0])
			}
		case dns.TypeNS:
			if nsRecord, ok := ans.(*dns.NS); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(nsRecord.Hdr.Ttl), nsRecord.Ns)
			}
		case dns.TypeSOA:
			if soaRecord, ok := ans.(*dns.SOA); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(soaRecord.Hdr.Ttl), soaRecord.Ns, soaRecord.Mbox)
			}
		case dns.TypePTR:
			if ptrRecord, ok := ans.(*dns.PTR); ok {
				fmt.Printf("%s\t%s.\t%s\t%s\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), FormatTTL(ptrRecord.Hdr.Ttl), ptrRecord.Ptr)
			}
		}
	}
}

func FormatTTL(ttl uint32) string {
	duration := time.Duration(ttl) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	if hours > 0 {
		return fmt.Sprintf("%02dh%02dm%02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%02dm%02ds   ", minutes, seconds)
	} else {
		return fmt.Sprintf("%02ds      ", seconds)
	}
}
