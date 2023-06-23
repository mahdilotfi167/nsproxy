/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package server

import (
	"log"
	"net"
	"nsproxy/config"
)

func RunServer(addr string, config *config.ServerConfig) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP port: %v", err)
	}
	defer conn.Close()

	log.Println("DNS Proxy server started on", addr)

	buffer := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		go handleDNSRequest(conn, addr, buffer[:n])
	}
}

func handleDNSRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	// Forward the DNS request to an upstream DNS server
	upstreamAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		log.Printf("Failed to resolve upstream DNS server address: %v\n", err)
		return
	}

	upstreamConn, err := net.DialUDP("udp", nil, upstreamAddr)
	if err != nil {
		log.Printf("Failed to establish connection to upstream DNS server: %v\n", err)
		return
	}
	defer upstreamConn.Close()

	// Send the DNS request to the upstream DNS server
	_, err = upstreamConn.Write(request)
	if err != nil {
		log.Printf("Failed to send DNS request to upstream DNS server: %v\n", err)
		return
	}

	// Receive the DNS response from the upstream DNS server
	responseBuf := make([]byte, 512)
	n, err := upstreamConn.Read(responseBuf)
	if err != nil {
		log.Printf("Failed to receive DNS response from upstream DNS server: %v\n", err)
		return
	}

	// Send the DNS response back to the client
	_, err = conn.WriteToUDP(responseBuf[:n], addr)
	if err != nil {
		log.Printf("Failed to send DNS response to client: %v\n", err)
		return
	}
}
