package randomid

import (
	"crypto/rand"
	"encoding/hex"
)

// Generator creates RFC 4122 version 4 identifiers from the operating system CSPRNG.
// A CSPRNG failure is unrecoverable because continuing could break durable identity.
type Generator struct{}

func (Generator) New() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic(err)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	encoded := hex.EncodeToString(value[:])
	return encoded[0:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:32]
}
