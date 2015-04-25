package main

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/gob"
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
	Version			int16
	Count			int16
	Sys_uptime		int32
	Unix_secs		int32
	Unix_nsecs		int32
	Flow_sequence		int32
	Engine_type		int8 
	Engine_id		int8 
	Sampling_interval	int16 
}

type NFLOW_v5_body struct {
	Srcaddr		int32
	Dstaddr		int32
	Nexthop		int32
	Input		int16
	Output		int16
	DPkts		int32
	DOctets		int32
	First		int32
	Last		int32
	Srcport		int16
	Dstport		int16
	Pad1		int8
	Tcp_flags	int8
	Prot		int8
	Tos		int8
	Src_as  	int16
	Dst_as  	int16
	Src_mask 	int8 
	Dst_mask	int8
	Pad2		int16
}

func construct_payload() gopacket.Payload {

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
	body := NFLOW_v5_body {
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
	}
	encBufHeader := new(bytes.Buffer)
	err := gob.NewEncoder(encBufHeader).Encode(header)
	if err != nil {
		log.Fatal(err)
	}

	encBufBody := new(bytes.Buffer)
	err = gob.NewEncoder(encBufBody).Encode(body)
	if err != nil {
		log.Fatal(err)
	}

	return gopacket.Payload(append(encBufHeader.Bytes(), encBufBody.Bytes()...))
}

func main() {

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	
	l2 := construct_ethernet()
	l3 := construct_ip("1.2.3.4", "5.6.7.8")
	l4 := construct_udp()
	payload := construct_payload()

	//LayerCake
	gopacket.SerializeLayers(buf, opts,
		l2, 
		l3,
		l4,
		payload)
//		gopacket.Payload([]byte{9, 10, 11, 12}))
	packetData := buf.Bytes()
        fmt.Println(packetData)
}
