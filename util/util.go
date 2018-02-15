package util

import (
	"crypto/rand"
	"io"
	"os"
	"strconv"
)

type UUID [16]byte

func GetEnvInt(key string) int {
	i, _ := strconv.Atoi(os.Getenv(key))
	return i
}

func GetEnvIntOrDefault(key string, def int) int {
	i, e := strconv.Atoi(os.Getenv(key))
	if e != nil {
		return def
	}
	return i
}

func GetEnvString(key string) string {
	return os.Getenv(key)
}

func GetEnvStringOrDefault(key string, def string) string {
	r := os.Getenv(key)
	if r == "" {
		return def
	}
	return r
}

func RandomUUIDString() string {
	u, err := randomUUID()

	if err != nil {
		panic(err)
	}

	return u.String()
}

func randomUUID() (UUID, error) {
	var u UUID
	_, err := io.ReadFull(rand.Reader, u[:])
	if err != nil {
		return u, err
	}
	u[6] &= 0x0F // clear version
	u[6] |= 0x40 // set version to 4 (random uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant
	return u, nil
}

func (u UUID) String() string {
	var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
	const hexString = "0123456789abcdef"
	r := make([]byte, 36)
	for i, b := range u {
		r[offsets[i]] = hexString[b>>4]
		r[offsets[i]+1] = hexString[b&0xF]
	}
	r[8] = '-'
	r[13] = '-'
	r[18] = '-'
	r[23] = '-'
	return string(r)

}
