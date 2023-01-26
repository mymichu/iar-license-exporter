package localhostserver

import (
	"log"
	"net"
)

type ServerSpecs struct {
	Hostname string
	Ipaddr   string
}

func GetServerSpecs() ServerSpecs {
	var hostname_list, _ = net.LookupHost("")
	hostname := "localhost"
	if len(hostname_list) > 0 {
		hostname = hostname_list[0]
	} else if len(hostname_list) > 1 {
		log.Fatal("More then 1 hostname found")
	}

	ipaddr, _ := net.LookupIP(hostname)
	return ServerSpecs{Hostname: hostname, Ipaddr: ipaddr[0].To4().String()}
}
