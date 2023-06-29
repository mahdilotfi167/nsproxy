/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package dns

const (
	A     uint16 = 1
	NS    uint16 = 2
	CNAME uint16 = 5
	MX    uint16 = 15
	TXT   uint16 = 16
	AAAA  uint16 = 28
)

const (
	IN uint16 = 1
	CH uint16 = 3
)

const ALL uint16 = 255

const (
	NO_ERROR        uint8 = 0
	FORMAT_ERROR    uint8 = 1
	SERVER_FAILURE  uint8 = 2
	NAME_ERROR      uint8 = 3
	NOT_IMPLEMENTED uint8 = 4
	REFUSED         uint8 = 5
)
