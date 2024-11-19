package validation

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
	ulidUtil "github.com/ucho456job/pocgo/pkg/ulid"
)

func ValidULID(ulid string) error {
	valid := ulidUtil.IsValid(ulid)
	if !valid {
		return ulidUtil.ErrInvalidULID
	}
	return nil
}

func ValidYYYYMMDD(yyyymmdd string) error {
	var yyyymmddRegex = regexp.MustCompile(`^\d{8}$`)
	return v.Validate(yyyymmdd, v.Match(yyyymmddRegex))
}

func ValidPage(page int) error {
	return v.Validate(page, v.Min(1))
}

func ValidSort(sort string) error {
	return v.Validate(sort, v.In("ASC", "DESC"))
}
