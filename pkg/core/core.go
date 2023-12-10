package core

import (
	"os"
	"runtime"
	"time"

	model "github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
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
func DisplayRecords(domainName string, answers []dns.RR) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Type", "Domain", "TTL", "Data"})
	var rows []table.Row
	for _, ans := range answers {
		var row table.Row
		switch ans.(type) {
		case *dns.A:
			aRecord, _ := ans.(*dns.A)
			row = table.Row{
				color.HiYellowString("A"),
				color.HiBlueString(domainName),
				aRecord.Hdr.Ttl,
				aRecord.A,
			}
		case *dns.AAAA:
			aaaaRecord, _ := ans.(*dns.AAAA)
			row = table.Row{
				color.HiYellowString("AAAA"),
				color.HiBlueString(domainName),
				aaaaRecord.Hdr.Ttl,
				aaaaRecord.AAAA,
			}
		case *dns.CNAME:
			cnameRecord, _ := ans.(*dns.CNAME)
			row = table.Row{
				"%s\t%s.\t%d\t%s\n",
				color.HiYellowString("CNAME"),
				color.HiBlueString(domainName),
				cnameRecord.Hdr.Ttl,
				cnameRecord.Target,
			}
		case *dns.MX:
			mxRecord, _ := ans.(*dns.MX)
			row = table.Row{
				color.HiYellowString("MX"),
				color.HiBlueString(domainName),
				mxRecord.Hdr.Ttl,
				mxRecord.Preference,
				mxRecord.Mx,
			}
		case *dns.TXT:
			txtRecord, _ := ans.(*dns.TXT)
			row = table.Row{
				color.HiYellowString("TXT"),
				color.HiBlueString(domainName),
				txtRecord.Hdr.Ttl,
				txtRecord.Txt[0],
			}
		case *dns.NS:
			nsRecord, _ := ans.(*dns.NS)
			row = table.Row{
				color.HiYellowString("NS"),
				color.HiBlueString(domainName),
				nsRecord.Hdr.Ttl,
				nsRecord.Ns,
			}
		case *dns.SOA:
			soaRecord, _ := ans.(*dns.SOA)
			row = table.Row{
				color.HiYellowString("SOA"),
				color.HiBlueString(domainName),
				soaRecord.Hdr.Ttl,
				soaRecord.Ns,
				soaRecord.Mbox,
			}
		case *dns.PTR:
			ptrRecord, _ := ans.(*dns.PTR)
			row = table.Row{
				color.HiYellowString("PTR"),
				color.HiBlueString(domainName),
				ptrRecord.Hdr.Ttl,
				ptrRecord.Ptr,
			}
		}
		rows = append(rows, row)
	}
	t.AppendRows(rows)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.Render()
}
