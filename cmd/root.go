package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bschaatsbergen/dnsee/pkg/core"
	"github.com/bschaatsbergen/dnsee/pkg/model"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var (
	version                string
	dnsServerIP            string
	userSpecifiedQueryType string

	rootCmd = &cobra.Command{
		Use:     "dnsee",
		Short:   "dnsee - Check DNS configurations quickly",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a domain name")
				fmt.Println("See 'dnsee -h' for help and examples")
				os.Exit(1)
			}

			domainName := args[0]

			client := dns.Client{}

			queryTypes := core.GetQueryTypes()

			// If a specific query type is provided, filter the queryTypes slice to only include that type
			if userSpecifiedQueryType != "" {
				queryTypes = filterQueryTypes(queryTypes, userSpecifiedQueryType)
			}

			for _, queryType := range queryTypes {
				msg := core.PrepareDNSQuery(domainName, queryType.Type)

				response, _, err := core.SendDNSQuery(&client, msg, dnsServerIP)
				if err != nil {
					log.Fatal(err)
				}

				core.DisplayRecords(domainName, queryType, response.Answer)
			}
		},
	}
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVar(&dnsServerIP, "dns-server-ip", "8.8.8.8", "IP address of the DNS server")
	rootCmd.Flags().StringVarP(&userSpecifiedQueryType, "query-type", "q", "", "Specific query type to filter on")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func filterQueryTypes(queryTypes []model.QueryType, userSpecifiedQueryType string) []model.QueryType {
	var filteredQueryTypes []model.QueryType
	for _, queryType := range queryTypes {
		if queryType.Name == userSpecifiedQueryType {
			filteredQueryTypes = append(filteredQueryTypes, queryType)
			break
		}
	}
	return filteredQueryTypes
}
