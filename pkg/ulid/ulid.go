package ulid

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"

	"github.com/oklog/ulid/v2"
	"github.com/ucho456job/pocgo/pkg/timer"
)

func New() string {
	return ulid.Make().String()
}

func IsValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

func GenerateStaticULID(seed string) string {
	hash := md5.Sum([]byte(seed))
	seedInt := binary.BigEndian.Uint64(hash[:8])
	source := rand.NewSource(int64(seedInt))
	entropy := rand.New(source)
	t := timer.Now()
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
