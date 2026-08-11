package shared

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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