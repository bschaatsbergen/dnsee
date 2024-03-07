package core

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	model "github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/fatih/color"
	"github.com/juju/ansiterm"
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

	logrus.Debugf("Sending DNS query to %s with query type: %s", dnsServerIP, dns.TypeToString[msg.Question[0].Qtype])
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

func DisplayRecords(domainName string, results []model.QueryResult) {

	w := ansiterm.NewTabWriter(os.Stdout, 8, 8, 4, ' ', 0)
	w.SetColorCapable(true)

	for _, result := range results {
		for _, record := range result.Records {
			switch result.QueryType.Type {
			case dns.TypeA:
				if aRecord, ok := record.(*dns.A); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(aRecord.Hdr.Ttl)), color.HiWhiteString(aRecord.A.String()))
				}
			case dns.TypeAAAA:
				if aaaaRecord, ok := record.(*dns.AAAA); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(aaaaRecord.Hdr.Ttl)), color.HiWhiteString(aaaaRecord.AAAA.String()))
				}
			case dns.TypeCNAME:
				if cnameRecord, ok := record.(*dns.CNAME); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(cnameRecord.Hdr.Ttl)), color.HiWhiteString(cnameRecord.Target))
				}
			case dns.TypeMX:
				if mxRecord, ok := record.(*dns.MX); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(mxRecord.Hdr.Ttl)), strings.Join([]string{color.HiRedString(strconv.FormatUint(uint64(mxRecord.Preference), 10)), color.HiWhiteString(mxRecord.Mx)}, "  "))
				}
			case dns.TypeTXT:
				if txtRecord, ok := record.(*dns.TXT); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(txtRecord.Hdr.Ttl)), color.HiWhiteString(txtRecord.Txt[0]))
				}
			case dns.TypeNS:
				if nsRecord, ok := record.(*dns.NS); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(nsRecord.Hdr.Ttl)), color.HiWhiteString(nsRecord.Ns))
				}
			case dns.TypeSOA:
				if soaRecord, ok := record.(*dns.SOA); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(soaRecord.Hdr.Ttl)), color.HiWhiteString(soaRecord.Ns), color.GreenString(soaRecord.Mbox))
				}
			case dns.TypePTR:
				if ptrRecord, ok := record.(*dns.PTR); ok {
					fmt.Fprintf(w, "%s\t%s.\t%s\t%s\t\n", color.HiYellowString(result.QueryType.Name), color.HiBlueString(domainName), color.HiMagentaString(FormatTTL(ptrRecord.Hdr.Ttl)), color.HiWhiteString(ptrRecord.Ptr))
				}
			}
		}
	}
	w.Flush() // Write table to stdout
}

func FormatTTL(ttl uint32) string {
	duration := time.Duration(ttl) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	if hours > 0 {
		return fmt.Sprintf("%02dh%02dm%02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%02dm%02ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%02ds", seconds)
	}
}
