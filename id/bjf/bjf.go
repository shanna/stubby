package bjf

import (
	"strconv"

	b "github.com/xor-gate/bjf"
)

// TODO: 2 byte checksum with secret just to increase the difficulty a notch.

type ID struct {
	secret []byte
}

func New(secret []byte) (*ID, error) {
	return &ID{secret}, nil
}

func (i ID) Encode(id string) string {
	return b.Encode(id)
}

func (i ID) Decode(token string) string {
	return strconv.Itoa(b.Decode(token))
}
