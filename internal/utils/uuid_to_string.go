package utils

import (
	"encoding/hex"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(uuid *pgtype.UUID) string {
	return string(encodeHex(uuid.Bytes))
}

func encodeHex(raw [16]byte) []byte {
	var result = make([]byte, 36)

	hex.Encode(result, raw[:4])
	result[8] = '-'
	hex.Encode(result[9:13], raw[4:6])
	result[13] = '-'
	hex.Encode(result[14:18], raw[6:8])
	result[18] = '-'
	hex.Encode(result[19:23], raw[8:10])
	result[23] = '-'
	hex.Encode(result[24:], raw[10:])

	return result
}
