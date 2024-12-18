package ulid

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var ErrInvalidULID = errors.New("invalid ulid")

func New() string {
	return ulid.Make().String()
}

func IsValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// GenerateStaticULID は引数の文字列に基づいて固定のULIDを生成します
// 同じ引数からは常に同じULIDが生成されます
func GenerateStaticULID(seed string) string {
	hash := md5.Sum([]byte(seed))
	seedInt := binary.BigEndian.Uint64(hash[:8])
	source := rand.NewSource(int64(seedInt))
	entropy := rand.New(source)
	fixedTime := ulid.Timestamp(time.Unix(0, 0))
	return ulid.MustNew(fixedTime, entropy).String()
}
