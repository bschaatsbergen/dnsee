package core

import (
	"fmt"
	"os"
	"runtime"
	"time"

	model "github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/fatih/color"
	"github.com/miekg/dns"
)

const Windows = "windows"
const ResolverPath = "/etc/resolv.conf"

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

func PrepareDNSQuery(domainName string, queryType uint16) dns.Msg {
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domainName), queryType)
	return msg
}

func SendDNSQuery(client *dns.Client, msg dns.Msg, dnsServerIP string) (*dns.Msg, time.Duration, error) {
	if dnsServerIP == "" {
		goOS := runtime.GOOS
		if goOS == Windows {
			fmt.Println("Unable to retrieve DNS configuration on Windows systems. \nPlease specify ip explicitely with the --dns-server-ip flag.")
			os.Exit(2)
		}
		conf, err := dns.ClientConfigFromFile(ResolverPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("Could not retrieve DNS server ip from system configuration. \nPlease specify ip explicitely with the --dns-server-ip flag.")
			os.Exit(2)
		}
		dnsServerIP = conf.Servers[0]
	}
	return client.Exchange(&msg, dnsServerIP+":53")
}

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
				fmt.Printf("%s\t%s.\t%d\t%s\t%d\n", color.HiYellowString(queryType.Name), color.HiBlueString(domainName), mxRecord.Hdr.Ttl, mxRecord.Mx, mxRecord.Preference)
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
