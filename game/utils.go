package gsnake

import (
	"encoding/hex"
	"math/rand"
)

func generateUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	// Set the version (4) and variant (RFC4122) bits
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	// Convert UUID to string format
	uuidStr := make([]byte, 36)
	hex.Encode(uuidStr[0:8], uuid[0:4])
	uuidStr[8] = '-'
	hex.Encode(uuidStr[9:13], uuid[4:6])
	uuidStr[13] = '-'
	hex.Encode(uuidStr[14:18], uuid[6:8])
	uuidStr[18] = '-'
	hex.Encode(uuidStr[19:23], uuid[8:10])
	uuidStr[23] = '-'
	hex.Encode(uuidStr[24:], uuid[10:])

	return string(uuidStr)
}
