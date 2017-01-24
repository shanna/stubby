package bjf

import (
	"hash/crc32"
	"strconv"

	b "github.com/xor-gate/bjf"
)

// TODO: 2 byte checksum with secret just to increase the difficulty a notch.
// TODO: ID interface needs error handling.

type ID struct {
	secret int
}

func New(secret []byte) (*ID, error) {
	return &ID{int(crc32.ChecksumIEEE(secret))}, nil
}

func (i ID) Encode(id string) string {
	in, _ := strconv.Atoi(id)
	return b.Encode(strconv.Itoa(i.secret ^ in))
}

func (i ID) Decode(token string) string {
	out := b.Decode(token)
	return strconv.Itoa(i.secret ^ out)
}
