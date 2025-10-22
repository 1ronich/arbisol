package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"arbisol/src/types"

	"github.com/mr-tron/base58"
)

func BigIntToLEBytesRightPad(bi *big.Int, totalSize int) []byte {
	// Get minimal big-endian byte slice
	if bi == nil {
		panic("BigIntToLEBytesRightPad: received nil *big.Int")
	}
	be := bi.Bytes()

	// Reverse it to little-endian
	le := make([]byte, len(be))
	for i := 0; i < len(be); i++ {
		le[i] = be[len(be)-1-i]
	}

	// Pad to totalSize by adding 0s to the right
	if len(le) < totalSize {
		padding := make([]byte, totalSize-len(le))
		le = append(le, padding...)
	}

	return le
}

func BigIntToLittleEndianBytes(n *big.Int) []byte {
	b := n.Bytes()
	// Pad with 0 if value is zero
	if len(b) == 0 {
		return []byte{0}
	}
	// Reverse for little-endian
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}
func ReadU8(buf *bytes.Reader) uint8 {
	var val uint8
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}
func ReadU16(buf *bytes.Reader) uint16 {
	var val uint16
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		panic(err)
	}
	return val
}
func ReadU32(buf *bytes.Reader) uint32 {
	var val uint32
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		panic(err)
	}
	return val
}
func ReadI32(buf *bytes.Reader) int32 {
	var val int32
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		panic(err)
	}
	return val
}
func ReadU64(buf *bytes.Reader) uint64 {
	var val uint64
	if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
		panic(err)
	}
	return val
}
func ReadU128(buf *bytes.Reader) types.Uint128 {
	var low, high uint64
	if err := binary.Read(buf, binary.LittleEndian, &low); err != nil {
		panic(err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &high); err != nil {
		panic(err)
	}
	return types.Uint128{Low: low, High: high}
}

func ReadPubkey(buf *bytes.Reader) string {
	key := make([]byte, 32)
	buf.Read(key)
	return base58.Encode(key)
}
