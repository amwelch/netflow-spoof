package main

import (
	"fmt"
	"net"
  "time"
	"encoding/binary"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
	"flag"
)

var NETFLOW_V5_HEADER_SIZE int = 24;
var NETFLOW_V5_RECORD_SIZE int = 48;
var PROTOCOL_TCP uint8 = 6
var PROTOCOL_UDP uint8 = 17
var NETFLOW_PORT int = 2055
var NANOSECOND int64 = 1000000000

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

type NETFLOW_v5_header struct {
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

type NETFLOW_v5_record struct {
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


func v4_to_uint32(addr net.IP) uint32 {
	var ret uint32;
	ret |= uint32(addr[0])
	ret |= uint32(addr[1]) << 8
	ret |= uint32(addr[2]) << 16
	ret |= uint32(addr[3]) << 24
	return ret
}

func construct_v5_header(count uint16, sampling uint16) NETFLOW_v5_header {
	header := NETFLOW_v5_header{
		Version:		5,		//Netflow v5
		Count:			count,		//Number of records in this packet
		Sys_uptime:		0,		//Ignore for now
		Unix_secs:		0,		//Ignore for now
		Unix_nsecs:		0,		//Ignore for now
		Flow_sequence:		0,		//Ignore for now. Eventually want to track sequence numbers
		Engine_type:		0,		//Ignore for now
		Engine_id:		0,		//Ignore for now
		Sampling_interval:	sampling,	//TODO
	}
	return header
}

func insert_v5_header(header NETFLOW_v5_header, buf []byte, offset int) int {
        binary.BigEndian.PutUint16(buf[offset:], header.Version)
        binary.BigEndian.PutUint16(buf[offset + 2:], header.Count)
        binary.BigEndian.PutUint32(buf[offset + 4:], header.Sys_uptime)
        binary.BigEndian.PutUint32(buf[offset + 8:], header.Unix_secs)
        binary.BigEndian.PutUint32(buf[offset + 12:], header.Unix_nsecs)
        binary.BigEndian.PutUint32(buf[offset + 16:], header.Flow_sequence)
        buf[offset + 20] = header.Engine_type
        buf[offset + 21] = header.Engine_id
        binary.BigEndian.PutUint16(buf[offset + 22:], header.Sampling_interval)

	return NETFLOW_V5_HEADER_SIZE
}

func construct_v5_record(srcaddr string, dstaddr string, 
	pkts uint32, l3_bytes uint32, srcport uint16, dstport uint16,
	protocol uint8, src_as uint16, dst_as uint16) NETFLOW_v5_record {

	srcip := v4_to_uint32(net.ParseIP(srcaddr))
	dstip := v4_to_uint32(net.ParseIP(dstaddr))

	record := NETFLOW_v5_record {
		Srcaddr:		srcip,
		Dstaddr:		dstip,
		Nexthop:		0,				//Ignore for now
		Input:			0,				//Do something with this later
		Output:			0,				//^^
		DPkts:			pkts,
		DOctets:		l3_bytes,
		First:			0,				//Ignore for now
		Last:			0,				//Ignore for now
		Srcport:		srcport,
		Dstport:		dstport,
		Pad1:			0,	
		Tcp_flags:		0,				//Something with this later
		Prot:			PROTOCOL_TCP,
		Tos:			0,
		Src_as:			src_as,
		Dst_as:			dst_as,
		Src_mask:		0,
		Dst_mask:		0,
	}
	return record
}

func insert_v5_record(record NETFLOW_v5_record, buf []byte, offset int) int {
        binary.BigEndian.PutUint32(buf[offset:], record.Srcaddr)
        binary.BigEndian.PutUint32(buf[offset + 4:], record.Dstaddr)
        binary.BigEndian.PutUint32(buf[offset + 8:], record.Nexthop)
        binary.BigEndian.PutUint16(buf[offset + 12:], record.Input)
        binary.BigEndian.PutUint16(buf[offset + 14:], record.Output)
        binary.BigEndian.PutUint32(buf[offset + 16:], record.DPkts)
        binary.BigEndian.PutUint32(buf[offset + 20:], record.DOctets)
        binary.BigEndian.PutUint32(buf[offset + 24:], record.First)
        binary.BigEndian.PutUint32(buf[offset + 28:], record.Last)
        binary.BigEndian.PutUint16(buf[offset + 32:], record.Srcport)
        binary.BigEndian.PutUint16(buf[offset + 34:], record.Dstport)
        buf[offset + 36] = record.Pad1
        buf[offset + 37] = record.Tcp_flags
        buf[offset + 38] = record.Prot
        buf[offset + 39] = record.Tos
        binary.BigEndian.PutUint16(buf[offset + 40:], record.Src_as)
        binary.BigEndian.PutUint16(buf[offset + 42:], record.Dst_as)
        buf[offset + 44] = record.Src_mask
        buf[offset + 45] = record.Dst_mask
        binary.BigEndian.PutUint16(buf[offset + 46:], record.Pad2)
	return NETFLOW_V5_RECORD_SIZE;
}

func construct_payload(num_records uint16) gopacket.Payload {

	buf := gopacket.NewSerializeBuffer()
//        payload := buf.Bytes()
        //Allocate the space we will need for the header
        bytes,err := buf.PrependBytes(NETFLOW_V5_HEADER_SIZE + NETFLOW_V5_RECORD_SIZE*int(num_records))
	if err != nil {
		return nil
	} 

	offset := 0

	header := construct_v5_header(num_records, 1000)
	offset += insert_v5_header(header, bytes, offset)

	var record NETFLOW_v5_record;
	for i := 0; i < int(num_records); i++ {
		record = construct_v5_record("1.1.1.1", "2.2.2.2", 5, 256, 80, 5050, 6, 237, 237)
		insert_v5_record(record, bytes, offset)
		
	}
        
	return gopacket.Payload(bytes)
}

func chk(err error){
	if err != nil {
		panic(err)
	}
}

func init_connection(addr net.IP) *net.UDPConn {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: addr, Port: 0})
	chk(err)
	return conn
}

func send_packet(conn *net.UDPConn, addr net.IP, port int, pkt []byte) {
	_, err := conn.WriteToUDP(pkt, &net.UDPAddr{IP: addr, Port: port})
	chk(err)
}


func main() {

	dst_ip := flag.String("dst", "127.0.0.1", "Destination IP to send the spoofed netflow")
	dst_port := flag.Int("port", NETFLOW_PORT, "Destination Port to send the spoofed netflow")
	rate := flag.Int("rate", 1, "Rate in MB/s")
  runtime := flag.Int64("time", 10, "Time in seconds to send packets")
  flows_per_packet := flag.Uint("fpp", 1, "flows per packet, max of 30")
  flag.Parse()

	dst_addr := net.ParseIP(*dst_ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	
	l2 := construct_ethernet()
	l3 := construct_ip("1.2.3.4", "5.6.7.8")
	l4 := construct_udp()

  start := time.Now().Unix()
  var last_loop int64
  var count int
  var cur_packet int64
  var wait int64   

  for last_loop - start < (*runtime)*NANOSECOND {
    count = 0
    cur_packet = time.Now().Unix()
    for count < *rate {
			payload := construct_payload(uint16(*flows_per_packet))
			//LayerCake
			gopacket.SerializeLayers(buf, opts,
				l2, 
				l3,
				l4,
				payload)
			packetData := buf.Bytes()
      
			//Send the packet to lo
			conn := init_connection(dst_addr)
			send_packet(conn, dst_addr, *dst_port, packetData)
      count++
      wait = 1*NANOSECOND - (time.Now().Unix() - cur_packet)
      if wait > 0 {
      	time.Sleep(time.Duration(wait))
      }
    }
    last_loop = time.Now().Unix()
	}
	fmt.Println("fin")
}
