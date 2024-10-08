package transaction_domain

import "github.com/ucho456job/pocgo/pkg/ulid"

func ValidID(id string) error {
	if !ulid.IsValid(id) {
		return ErrInvalidTransactionID
	}
	return nil
}
