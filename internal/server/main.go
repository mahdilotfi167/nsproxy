/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package server

import (
	"context"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"log"
	"net"
	"nsproxy/config"
	"nsproxy/pkg/dns"
	"time"
)

type Server struct {
	addr   string
	config *config.ServerConfig
	cache  *cache.Cache[string]
}

func NewServer(addr string, config *config.ServerConfig, cache *cache.Cache[string]) *Server {
	return &Server{
		addr:   addr,
		config: config,
		cache:  cache,
	}
}

func (s *Server) Run() {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP port: %v", err)
	}
	defer conn.Close()

	log.Println("DNS Proxy server started on", s.addr)

	buffer := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		ctx := context.Background()
		go s.handleDNSRequest(ctx, conn, addr, buffer[:n])
	}
}

func (s *Server) handleDNSRequest(ctx context.Context, conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	var err error
	// Create required DNS messages
	requestMsg, err := dns.ParseMessage(request)

	extRequestMsg := dns.Message{
		ID: requestMsg.ID,
	}

	responseMsg := dns.Message{
		ID:        requestMsg.ID,
		Flags:     dns.Flags{QR: true, Opcode: requestMsg.Flags.Opcode},
		Questions: requestMsg.Questions,
	}

	// Separate existing DNS questions versus missing questions
	for _, question := range requestMsg.Questions {
		qBytes := QuestionToBytes(question)
		if value, err := s.cache.Get(ctx, string(qBytes)); err == nil {
			resourceRecords := BytesToResourceRecords([]byte(value))
			responseMsg.Answers = append(responseMsg.Answers, resourceRecords...)
		} else {
			extRequestMsg.Questions = append(extRequestMsg.Questions, question)
		}
	}

	// Resolve missing DNS questions via external DNS server
	if len(extRequestMsg.Questions) > 0 {
		timeout := s.config.ExternalDNSTimeout
		if timeout == 0 {
			timeout = 60
		}

		extResponseMsg, err := resolveDNSRequest(*requestMsg, s.config.ExternalDNSServers, time.Duration(timeout)*time.Second)

		if err != nil {
			responseMsg.Flags.ResponseCode = dns.REFUSED
		} else {
			s.putCache(ctx, extResponseMsg)

			responseMsg.Answers = append(responseMsg.Answers, extResponseMsg.Answers...)
		}
	}

	// Send the DNS response back to the client
	_, err = conn.WriteToUDP(responseMsg.ToBytes(), addr)
	if err != nil {
		log.Printf("Failed to send DNS response to client: %v\n", err)
		return
	}
}

func (s *Server) putCache(ctx context.Context, msg *dns.Message) {
	records := make([]dns.ResourceRecord, len(msg.Answers)+len(msg.AuthorityRecords)+len(msg.AdditionalRecords))
	records = append(records, msg.Answers...)
	records = append(records, msg.AuthorityRecords...)
	records = append(records, msg.AdditionalRecords...)

	questionToRR := make(map[dns.Question][]dns.ResourceRecord)
	for _, record := range records {
		question := dns.Question{Name: record.Name, Type: record.Type, Class: record.Class}
		questionToRR[question] = append(questionToRR[question], record)
	}

	for question, records := range questionToRR {
		qBytes := QuestionToBytes(question)
		rrBytes := ResourceRecordsToBytes(records)
		s.cache.Set(ctx, string(qBytes), string(rrBytes), store.WithExpiration(time.Duration(s.config.CacheExpirationTime)*time.Second))
	}
}

func resolveDNSRequest(dnsMsg dns.Message, dnsServers []string, timeout time.Duration) (*dns.Message, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Create a channel to receive the resolved DNS message
	resultCh := make(chan *dns.Message, 1)

	// Resolve DNS asynchronously for each DNS server in the list
	for _, dnsServer := range dnsServers {
		go func(server string) {
			// Create a UDP connection to the DNS server
			conn, err := net.Dial("udp", server)
			if err != nil {
				return
			}
			defer conn.Close()

			if timeout > 0 {
				conn.SetDeadline(time.Now().Add(timeout))
			}

			dnsBytes := dnsMsg.ToBytes()

			_, err = conn.Write(dnsBytes)
			if err != nil {
				return
			}

			responseBuf := make([]byte, 512)
			n, err := conn.Read(responseBuf)
			if err != nil {
				return
			}

			// Parse the DNS response
			responseMsg, err := dns.ParseMessage(responseBuf[:n])
			if err != nil {
				return
			}

			// Check the response code
			if responseMsg.Flags.ResponseCode != dns.NO_ERROR {
				return
			}

			// Send the resolved DNS message to the result channel
			resultCh <- responseMsg
		}(dnsServer)
	}

	// Wait for the first resolved DNS message or error
	select {
	case response := <-resultCh:
		return response, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("DNS resolution timed out after %s", timeout.String())
	}
}
