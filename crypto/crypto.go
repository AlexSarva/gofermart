package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/google/uuid"
)

// SecretKey - secret key of encrypt and decrypt
var SecretKey = []byte("Ag@th@")

// ErrNotValidSing - error that occurs when it is impossible to recover UserID
var ErrNotValidSing = errors.New("sign is not valid")

// Encrypt convert user uuid to hash
// secret key should be the same for Encrypt and Decrypt
func Encrypt(uuid uuid.UUID, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(uuid[:])
	dst := h.Sum(nil)
	var fullCookie []byte
	fullCookie = append(fullCookie, uuid[:]...)
	fullCookie = append(fullCookie, dst...)
	return hex.EncodeToString(fullCookie)
}

// Decrypt convert user hash to uuid
func Decrypt(hashString string, secret []byte) (uuid.UUID, error) {
	var (
		data []byte // декодированное сообщение с подписью
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(hashString)
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, ErrNotValidSing
	}
	id, idErr := uuid.FromBytes(data[:16])
	if idErr != nil {
		log.Println(idErr)
	}
	h := hmac.New(sha256.New, secret)
	h.Write(data[:16])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[16:]) {
		return id, nil
	} else {
		return uuid.UUID{}, ErrNotValidSing
	}
}
