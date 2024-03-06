package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bschaatsbergen/dnsee/pkg/core"
	"github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/fatih/color"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type PlainFormatter struct{}

var (
	version   string
	flagStore model.Flagstore

	//nolint:goconst
	rootCmd = &cobra.Command{
		Use:     "dnsee",
		Short:   "dnsee - check DNS configurations quickly",
		Version: version,
		PreRun:  toggleDebug,
		Example: "dnsee " + color.New(color.FgBlue).SprintFunc()("example.com") + "." +
			"\n" + "dnsee " + color.New(color.FgBlue).SprintFunc()("example.com") + "." + " -q A" +
			"\n" + "dnsee " + color.New(color.FgBlue).SprintFunc()("example.com") + "." + " --dns-server-ip 8.8.8.8" +
			"\n" + "dnsee " + color.New(color.FgBlue).SprintFunc()("example.com") + "." + " --debug",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a domain name")
				fmt.Println("See 'dnsee -h' for help and examples")
				os.Exit(1)
			}

			domainName := args[0]

			client := dns.Client{}

			supportedQueryTypes := core.GetQueryTypes()

			typesToQuery := supportedQueryTypes

			// If a specific query type is provided, filter the queryTypes slice to only include that type
			if len(flagStore.UserSpecifiedQueryTypes) > 0 {
				typesToQuery = core.FilterQueryTypes(supportedQueryTypes, flagStore.UserSpecifiedQueryTypes)
			}

			// Send a DNS query for each query type in the queryTypes slice
			for _, queryType := range typesToQuery {
				msg := core.PrepareDNSQuery(domainName, queryType.Type)

				response, _, err := core.SendDNSQuery(&client, msg, flagStore.DNSServerIP)
				if err != nil {
					log.Fatal(err)
				}

				core.DisplayRecords(domainName, queryType, response.Answer)
			}
		},
	}
)

func init() {
	setupCobraUsageTemplate()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVar(&flagStore.DNSServerIP, "dns-server-ip", "", "IP address of the DNS server")
	rootCmd.Flags().StringSliceVarP(&flagStore.UserSpecifiedQueryTypes, "query-types", "q", queryTypesToStrings(core.GetQueryTypes()), "specific query type(s) to filter on")
	rootCmd.Flags().BoolVarP(&flagStore.Debug, "debug", "d", false, "verbose logging")
}

func setupCobraUsageTemplate() {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func queryTypesToStrings(queryTypes []model.QueryType) []string {
	var strings []string
	for _, s := range queryTypes {
		strings = append(strings, s.Name)
	}
	return strings
}

func toggleDebug(cmd *cobra.Command, args []string) {
	if flagStore.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logs enabled")
		log.SetFormatter(&log.TextFormatter{})
	} else {
		plainFormatter := new(PlainFormatter)
		log.SetFormatter(plainFormatter)
	}
}
