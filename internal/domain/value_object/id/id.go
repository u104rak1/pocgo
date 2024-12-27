package id

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidULID = errors.New("invalid ulid")
	ErrEmptyID     = errors.New("id must not be empty")
)

// ID ジェネリックなID型
type ID[T any] struct {
	value string
}

func New[T any]() ID[T] {
	return ID[T]{value: ulid.Make().String()}
}

func NewFromString[T any](value string) (ID[T], error) {
	if value == "" {
		return ID[T]{}, ErrEmptyID
	}
	if !IsValid(value) {
		return ID[T]{}, ErrInvalidULID
	}
	return ID[T]{value: value}, nil
}

func (id ID[T]) String() string {
	return id.value
}

func (id ID[T]) Equals(other ID[T]) bool {
	return id.value == other.value
}

func (id ID[T]) IsValid() bool {
	return id.value != "" && IsValid(id.value)
}

func NewForTest[T any](seed string) ID[T] {
	return ID[T]{value: GenerateStaticULID(seed)}
}

func IsValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// GenerateStaticULID は引数の文字列に基づいて固定のULIDを生成します
// 同じ引数からは常に同じULIDが生成されます
func GenerateStaticULID(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	seedInt := binary.BigEndian.Uint64(hash[:8])

	// uint64からint64への変換でオーバーフローの可能性を示唆していますが、テストで使用する為だけの関数な為、オーバーフローは考慮しません。
	// #nosec: G115
	source := rand.NewSource(int64(seedInt))

	// 　math/randは暗号学的に安全ではない乱数生成器なのでセキュリティが重要な場合は crypto/rand を使用すべきですが、テストで使用する為だけの関数な為、セキュリティは考慮しません。
	// #nosec: G404
	entropy := rand.New(source)

	fixedTime := ulid.Timestamp(time.Unix(0, 0))
	return ulid.MustNew(fixedTime, entropy).String()
}
