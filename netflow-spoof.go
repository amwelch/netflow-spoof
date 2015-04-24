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

type NFLOW_v5_header struct {
	version			int16
	count			int16
	sys_uptime		int32
	unix_secs		int32
	unix_nsecs		int32
	flow_sequence		int32
	engine_type		int8 
	engine_id		int8 
	sampling_interval	int16 
}

type NFLOW_v5_body struct {
	srcaddr		int32
	dstaddr		int32
	nexthop		int32
	input		int16
	output		int16
	dPkts		int32
	dOctets		int32
	first		int32
	last		int32
	srcport		int16
	dstport		int16
	pad1		int8
	tcp_flags	int8
	prot		int8
	tos		int8
	src_as  	int16
	dst_as  	int16
	src_mask 	int8 
	dst_mask	int8
	pad2		int16
}

func construct_nflow_v5() {

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
