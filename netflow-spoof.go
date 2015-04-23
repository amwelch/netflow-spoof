package main

import (
	"fmt"
	"net"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
)

func main() {

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	
	l2 := &layers.Ethernet{}
	l3 := &layers.IPv4{
		SrcIP: net.IP{1, 2, 3, 4},
		DstIP: net.IP{5, 6, 7, 8},
        }
	l4 := &layers.TCP{}


	//LayerCake
	gopacket.SerializeLayers(buf, opts,
		l2, 
		l3,
		l4,
		gopacket.Payload([]byte{1, 2, 3, 4}))
	packetData := buf.Bytes()
        fmt.Println(packetData)
}
