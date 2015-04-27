package main

import (
	"fmt"
	"net"
        "encoding/binary"
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
	Version			uint16
	Count			uint16
	Sys_uptime		uint32
	Unix_secs		uint32
	Unix_nsecs		uint32
	Flow_sequence		uint32
	Engine_type		uint8 
	Engine_id		uint8 
	Sampling_interval	uint16 
}

type NFLOW_v5_body struct {
	Srcaddr		uint32
	Dstaddr		uint32
	Nexthop		uint32
	Input		uint16
	Output		uint16
	DPkts		uint32
	DOctets		uint32
	First		uint32
	Last		uint32
	Srcport		uint16
	Dstport		uint16
	Pad1		uint8
	Tcp_flags	uint8
	Prot		uint8
	Tos		uint8
	Src_as  	uint16
	Dst_as  	uint16
	Src_mask 	uint8 
	Dst_mask	uint8
	Pad2		uint16
}


//func construct_payload() gopacket.Payload {
func construct_payload() {

	header := NFLOW_v5_header{
		Version:		0,
		Count:			0,
		Sys_uptime:		0,
		Unix_secs:		0,
		Unix_nsecs:		0,
		Flow_sequence:		0,
		Engine_type:		0,
		Engine_id:		0,
		Sampling_interval:	0,
	}
/*	body := NFLOW_v5_body {
		Srcaddr:		0,
		Dstaddr:		0,
		Nexthop:		0,
		Input:			0,
		Output:			0,
		DPkts:			0,
		DOctets:		0,
		First:			0,
		Last:			0,
		Srcport:		0,
		Dstport:		0,
		Pad1:			0,	
		Tcp_flags:		0,
		Prot:			0,
		Tos:			0,
		Src_as:			0,
		Dst_as:			0,
		Src_mask:		0,
		Dst_mask:		0,
	}*/

        NETFLOW_V5_HEADER_SIZE := 24;
 //       NETFLOW_V5_BODY_SIZE := 48;

	buf := gopacket.NewSerializeBuffer()
//        payload := buf.Bytes()
        //Allocate the space we will need for the header
        bytes,err := buf.PrependBytes(NETFLOW_V5_HEADER_SIZE)
	if err != nil {
		return 
	} 

        //Go through and add each field to the buffer
        binary.BigEndian.PutUint16(bytes, header.Version)
        
        
        fmt.Println(fmt.Sprintf("%v", bytes))

	return
//	return gopacket.Payload(buf)
}

func main() {

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	
	l2 := construct_ethernet()
	l3 := construct_ip("1.2.3.4", "5.6.7.8")
	l4 := construct_udp()
	construct_payload()

	//LayerCake
	gopacket.SerializeLayers(buf, opts,
		l2, 
		l3,
		l4,
		)
//		gopacket.Payload([]byte{9, 10, 11, 12}))
	packetData := buf.Bytes()
        fmt.Println("Entire Packet")
        fmt.Println(packetData)
}
