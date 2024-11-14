package validation

import (
	ulidUtil "github.com/ucho456job/pocgo/pkg/ulid"
)

func ValidULID(ulid string) error {
	valid := ulidUtil.IsValid(ulid)
	if !valid {
		return ulidUtil.ErrInvalidULID
	}
	return nil
}
