package validation

import (
	"errors"
	"regexp"
	"time"

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
	if err := v.Validate(yyyymmdd, v.Match(yyyymmddRegex)); err != nil {
		return err
	}

	_, err := time.Parse("20060102", yyyymmdd)
	if err != nil {
		return err
	}

	return nil
}

func ValidateDateRange(from, to string) error {
	fromDate, err := time.Parse("20060102", from)
	if err != nil {
		return err
	}
	toDate, err := time.Parse("20060102", to)
	if err != nil {
		return err
	}
	if toDate.Before(fromDate) {
		return errors.New("to date cannot be before from date")
	}
	return nil
}

func ValidPage(page int) error {
	if page <= 0 {
		return errors.New("page must be greater than 0")
	}
	return nil
}

func ValidSort(sort string) error {
	return v.Validate(sort, v.In("ASC", "DESC"))
}
