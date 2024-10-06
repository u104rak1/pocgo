package password

import (
	"encoding/base64"
	"errors"
)

var ErrPasswordUnmatch = errors.New("password unmatch")

func Encode(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func Compare(hash, password string) error {
	decodedPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	if string(decodedPassword) != password {
		return ErrPasswordUnmatch
	}
	return nil
}
