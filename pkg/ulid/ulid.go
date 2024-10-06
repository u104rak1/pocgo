// パッケージ ulid は、ULID (Universally Unique Lexicographically Sortable Identifiers) を生成および検証するためのユーティリティ関数を提供します。
package ulid

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULIDを生成する
func New() string {
	return ulid.Make().String()
}

// ULIDを検証する
func IsValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// 引数に指定した文字列を元にULIDを生成する
func GenerateStaticULID(seed string) string {
	hash := md5.Sum([]byte(seed))
	seedInt := binary.BigEndian.Uint64(hash[:8])
	source := rand.NewSource(int64(seedInt))
	entropy := rand.New(source)
	t := time.Now().UTC()
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
