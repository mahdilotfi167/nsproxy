package dns

type Message struct {
	ID                uint16
	Flags             Flags
	Questions         []Question
	Answers           []ResourceRecord
	AuthorityRecords  []ResourceRecord
	AdditionalRecords []ResourceRecord
}

type Flags struct {
	QR           bool  // Is Response
	Opcode       uint8 // Operation Code
	AA           bool  // Authoritative Answer
	TC           bool  // Truncated Response
	RD           bool  // Recursion Desired
	RA           bool  // Recursion Available
	Z            uint8 // Reserved for future use (should be zero)
	ResponseCode uint8 // Response Code
}

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

type ResourceRecord struct {
	Name    string
	Type    uint16
	Class   uint16
	TTL     uint32
	DataLen uint16
	Data    []byte
}
