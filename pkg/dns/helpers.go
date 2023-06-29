package dns

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func ParseMessage(data []byte) (*Message, error) {
	var msg Message

	// Parse the DNS header
	msg.ID = binary.BigEndian.Uint16(data[0:2])
	flags := binary.BigEndian.Uint16(data[2:4])
	msg.Flags = parseFlags(flags)

	// Parse the DNS questions
	questionCount := binary.BigEndian.Uint16(data[4:6])
	questionIndex := 12
	for i := 0; i < int(questionCount); i++ {
		question, newIndex := parseQuestion(data, questionIndex)
		msg.Questions = append(msg.Questions, question)
		questionIndex = newIndex
	}

	// Parse the DNS resource records
	answerCount := binary.BigEndian.Uint16(data[6:8])
	authorityCount := binary.BigEndian.Uint16(data[8:10])
	additionalCount := binary.BigEndian.Uint16(data[10:12])

	// Parse the answer resource records
	answerIndex := questionIndex
	for i := 0; i < int(answerCount); i++ {
		answer, newIndex := parseResourceRecord(data, answerIndex)
		msg.Answers = append(msg.Answers, answer)
		answerIndex = newIndex
	}

	// Parse the authority resource records
	authorityIndex := answerIndex
	for i := 0; i < int(authorityCount); i++ {
		authority, newIndex := parseResourceRecord(data, authorityIndex)
		msg.AuthorityRecords = append(msg.AuthorityRecords, authority)
		authorityIndex = newIndex
	}

	// Parse the additional resource records
	additionalIndex := authorityIndex
	for i := 0; i < int(additionalCount); i++ {
		additional, newIndex := parseResourceRecord(data, additionalIndex)
		msg.AdditionalRecords = append(msg.AdditionalRecords, additional)
		additionalIndex = newIndex
	}

	return &msg, nil
}

func parseFlags(flags uint16) Flags {
	var dnsFlags Flags

	dnsFlags.ResponseCode = uint8(flags & ((1 << 4) - 1))
	dnsFlags.Z = uint8((flags & ((1 << 7) - 1)) >> 4)
	dnsFlags.RA = flags&(1<<7) != 0
	dnsFlags.RD = flags&(1<<8) != 0
	dnsFlags.TC = flags&(1<<9) != 0
	dnsFlags.AA = flags&(1<<10) != 0
	dnsFlags.Opcode = uint8((flags & ((1 << 15) - 1)) >> 11)
	dnsFlags.QR = flags&(1<<15) != 0

	return dnsFlags
}

func parseQuestion(data []byte, index int) (Question, int) {
	var question Question

	// Parse the domain name
	domainName, newIndex := parseDomainName(data, index)
	question.Name = domainName

	// Parse the type and class
	question.Type = binary.BigEndian.Uint16(data[newIndex : newIndex+2])
	question.Class = binary.BigEndian.Uint16(data[newIndex+2 : newIndex+4])

	newIndex += 4

	return question, newIndex
}

func parseResourceRecord(data []byte, index int) (ResourceRecord, int) {
	var record ResourceRecord

	// Parse the domain name
	domainName, newIndex := parseDomainName(data, index)
	record.Name = domainName

	// Parse the type and class
	record.Type = binary.BigEndian.Uint16(data[newIndex : newIndex+2])
	record.Class = binary.BigEndian.Uint16(data[newIndex+2 : newIndex+4])

	// Parse the TTL
	record.TTL = binary.BigEndian.Uint32(data[newIndex+4 : newIndex+8])

	// Parse the data length
	dataLen := binary.BigEndian.Uint16(data[newIndex+8 : newIndex+10])
	record.DataLen = dataLen

	// Parse the data
	record.Data = data[newIndex+10 : newIndex+10+int(dataLen)]

	newIndex = newIndex + 10 + int(dataLen)

	return record, newIndex
}

func parseDomainName(data []byte, index int) (string, int) {
	var domainName string
	var currentIndex = index

	for {
		labelLen := int(data[currentIndex])

		if labelLen == 0 {
			// End of domain name
			break
		}

		if labelLen&0xC0 == 0xC0 {
			// Compressed label
			pointer := binary.BigEndian.Uint16(data[currentIndex : currentIndex+2])
			pointer &= 0x3FFF // Remove compression flag

			// Parse compressed label recursively
			compressedName, _ := parseDomainName(data, int(pointer))
			domainName += compressedName

			currentIndex += 1 // Move to next position
			break
		}

		label := string(data[currentIndex+1 : currentIndex+1+labelLen])
		domainName += label + "."

		currentIndex += labelLen + 1 // Move to next label
	}

	return domainName, currentIndex + 1
}

func (msg *Message) ToBytes() []byte {
	var buf bytes.Buffer

	// Write the ID field
	binary.Write(&buf, binary.BigEndian, msg.ID)

	// Write the flags field
	flags := msg.Flags.ToUint16()
	binary.Write(&buf, binary.BigEndian, flags)

	// Write the questions count
	binary.Write(&buf, binary.BigEndian, uint16(len(msg.Questions)))

	// Write the answers count
	binary.Write(&buf, binary.BigEndian, uint16(len(msg.Answers)))

	// Write the authority RRs count
	binary.Write(&buf, binary.BigEndian, uint16(len(msg.AuthorityRecords)))

	// Write the additional RRs count
	binary.Write(&buf, binary.BigEndian, uint16(len(msg.AdditionalRecords)))

	// Write the questions
	for _, question := range msg.Questions {
		buf.Write(question.ToBytes())
	}

	//Write the answers, authority RRs, and additional RRs
	for _, rr := range msg.Answers {
		buf.Write(rr.ToBytes())
	}
	for _, rr := range msg.AuthorityRecords {
		buf.Write(rr.ToBytes())
	}
	for _, rr := range msg.AdditionalRecords {
		buf.Write(rr.ToBytes())
	}

	return buf.Bytes()
}

func (f Flags) ToUint16() uint16 {
	var value uint16

	value |= uint16(f.ResponseCode)
	value |= uint16(f.Z) << 4
	if f.RA {
		value |= 1 << 7
	}
	if f.RD {
		value |= 1 << 8
	}
	if f.TC {
		value |= 1 << 9
	}
	if f.AA {
		value |= 1 << 10
	}
	value |= uint16(f.Opcode) << 11
	if f.QR {
		value |= 1 << 15
	}

	return value
}

func (q Question) ToBytes() []byte {
	var buf bytes.Buffer

	// Write the domain name
	buf.Write(encodeDomainName(q.Name))

	// Write the question type
	binary.Write(&buf, binary.BigEndian, q.Type)

	// Write the question class
	binary.Write(&buf, binary.BigEndian, q.Class)

	return buf.Bytes()
}

func (rr ResourceRecord) ToBytes() []byte {
	var buf bytes.Buffer

	// Write the domain name
	buf.Write(encodeDomainName(rr.Name))

	// Write the record type
	binary.Write(&buf, binary.BigEndian, rr.Type)

	// Write the record class
	binary.Write(&buf, binary.BigEndian, rr.Class)

	// Write the record TTL
	binary.Write(&buf, binary.BigEndian, rr.TTL)

	// Write the record data length
	dataLength := uint16(len(rr.Data))
	binary.Write(&buf, binary.BigEndian, dataLength)

	// Write the record data
	buf.Write(rr.Data)

	return buf.Bytes()
}

func encodeDomainName(name string) []byte {
	var buf bytes.Buffer

	labels := strings.Split(strings.TrimSuffix(name, "."), ".")
	for _, label := range labels {
		buf.WriteByte(uint8(len(label)))
		buf.WriteString(label)
	}

	buf.WriteByte(0x00) // Null terminator

	return buf.Bytes()
}
