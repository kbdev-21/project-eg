package shared

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ForceParseStringToUuid(s string) pgtype.UUID {
	uid := uuid.MustParse(s)
	return pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}
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