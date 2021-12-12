package waiter

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"github.com/google/uuid"
)

func GenerateUuid() string {
	id := uuid.New()
	return base64.RawURLEncoding.EncodeToString(id[:])
}

func GetBytes(d interface{}) ([]byte, error) {
	if d == nil {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
