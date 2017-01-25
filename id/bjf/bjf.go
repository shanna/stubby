package bjf

import (
	"fmt"
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

func (i ID) Encode(id string) (string, error) {
	id32, err := strconv.Atoi(id)
	if err != nil {
		return "", fmt.Errorf("ID encode error %s", err)
	}
	return b.Encode(strconv.Itoa(i.secret ^ id32)), nil
}

func (i ID) Decode(token string) (string, error) {
	out := b.Decode(token)
	return strconv.Itoa(i.secret ^ out), nil
}
