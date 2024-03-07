package model

import "github.com/miekg/dns"

type QueryType struct {
	Type uint16
	Name string
}

type QueryResult struct {
	QueryType QueryType
	Records   []dns.RR
}
