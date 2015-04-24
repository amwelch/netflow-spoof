package main

import (
	"fmt"
	"net"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
)

func construct_ethernet() *layers.Ethernet {
	return &layers.Ethernet{}
}

func construct_ip(srcaddr string, dstaddr string) *layers.IPv4 {
	return &layers.IPv4{
		SrcIP: net.ParseIP(srcaddr),
		DstIP: net.ParseIP(dstaddr),
	}
}

func construct_udp() *layers.UDP {
	return &layers.UDP{}
}

func main() {

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	
	l2 := construct_ethernet()
	l3 := construct_ip("1.2.3.4", "5.6.7.8")
	l4 := construct_udp()


	//LayerCake
	gopacket.SerializeLayers(buf, opts,
		l2, 
		l3,
		l4,
		gopacket.Payload([]byte{9, 10, 11, 12}))
	packetData := buf.Bytes()
        fmt.Println(packetData)
}
