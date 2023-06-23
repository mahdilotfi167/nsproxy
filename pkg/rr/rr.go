/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package rr

const TYPE_A uint16 = 1
const TYPE_NS uint16 = 2
const TYPE_CNAME uint16 = 5
const TYPE_MX uint16 = 15
const TYPE_TXT uint16 = 16
const TYPE_AAAA uint16 = 28
const TYPE_ALL uint16 = 255

const CLASS_IN uint16 = 1
const CLASS_CH uint16 = 3
const CLASS_ALL uint16 = 255

type Record struct {
	Name  string
	Value string
	Type  uint16
	Class uint16
	TTL   int32
}
