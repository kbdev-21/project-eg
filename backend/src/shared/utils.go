package shared

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func ParseStringToUuid(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if(err != nil) {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}, nil
}

func CapitalizeString(s string) string {
	if s == "" {
		return ""
	}

	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

func RandomInt(from int, to int) int {
	return rand.Intn(to-from+1) + from
}

// NormalizedText chuẩn hoá text tiếng Việt: viết thường, bỏ dấu (á, à, ả, ã, ạ,
// â, ê, ô, ơ, ư...), đổi đ/Đ thành d, và gộp khoảng trắng thừa.
func NormalizedText(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "đ", "d")

	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, err := transform.String(t, s)
	if err != nil {
		result = s
	}

	return strings.Join(strings.Fields(result), " ")
}