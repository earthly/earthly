package stringutil

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//RandomAlphanumeric returns a random alphanumeric string of length n
func RandomAlphanumeric(n int) string {
	var seed int64
	err := binary.Read(crand.Reader, binary.BigEndian, &seed)
	if err != nil {
		panic(err) // unlikely (if this did happen, the system would be very sick)
	}
	rand := mrand.New(mrand.NewSource(seed))

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
