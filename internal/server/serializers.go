package server

import (
	"bytes"
	"nsproxy/pkg/dns"
)

// QuestionToBytes converts a Question struct to a byte array
func QuestionToBytes(q dns.Question) []byte {
	return q.ToBytes()
}

// ResourceRecordsToBytes converts an array of ResourceRecord structs to a byte array
func ResourceRecordsToBytes(records []dns.ResourceRecord) []byte {
	buf := new(bytes.Buffer)

	for _, record := range records {
		buf.Write(record.ToBytes())
	}

	return buf.Bytes()
}

// BytesToResourceRecords converts a byte array to an array of ResourceRecord structs
func BytesToResourceRecords(data []byte) []dns.ResourceRecord {
	return dns.ParseResourceRecords(data)
}
