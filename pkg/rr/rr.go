/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package rr

const (
	TYPE_A     uint16 = 1
	TYPE_NS    uint16 = 2
	TYPE_CNAME uint16 = 5
	TYPE_MX    uint16 = 15
	TYPE_TXT   uint16 = 16
	TYPE_AAAA  uint16 = 28
	TYPE_ALL   uint16 = 255
)

const (
	CLASS_IN  uint16 = 1
	CLASS_CH  uint16 = 3
	CLASS_ALL uint16 = 255
)

type Record struct {
	Name  string
	Value string
	Type  uint16
	Class uint16
	TTL   int32
}
