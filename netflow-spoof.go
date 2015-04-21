package main

import (
	"fmt"
	"net"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
)

func main() {
	ip := &layers.IPv4{
		SrcIP: net.IP{1, 2, 3, 4},
		DstIP: net.IP{5, 6, 7, 8},
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err := ip.SerializeTo(buf, opts)
	if err != nil {panic(err) }
	
	fmt.Println(buf.Bytes())
}
