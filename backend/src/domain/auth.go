package domain

import (
	"backend/src/db"
	"backend/src/shared"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	db.User
	Role    UserRole    `json:"role"`
	AvtCode UserAvtCode `json:"avtCode"`
}

type UserRole string

const (
	ADMIN UserRole = "ADMIN"
	USER  UserRole = "USER"
	GUEST UserRole = "GUEST"
)

type UserAvtCode string

const (
	BUNNY   UserAvtCode = "BUNNY"
	KITTEN  UserAvtCode = "KITTEN"
	GRIZZLE UserAvtCode = "GRIZZLE"
	HAMSTER UserAvtCode = "HAMSTER"
	MONKEY  UserAvtCode = "MONKEY"
)

func ToUser(u db.User) User {
	return User{
		User:    u,
		Role:    UserRole(u.Role),
		AvtCode: UserAvtCode(u.AvtCode),
	}
}

func GetUserById(dbe *DbExec, id pgtype.UUID) (User, error) {
	dbU, err := dbe.q.GetUserById(dbe.c, id)
	return ToUser(dbU), err
}

func SyncUserFromTokenPayload(dbe *DbExec, payload TokenPayload) (User, error) {
	idFromPayload, err := shared.ParseStringToUuid(payload.Sub)
	if err != nil {
		return User{}, err
	}

	updateRole := USER
	updateEmail := pgtype.Text{
		String: payload.Email,
		Valid:  true,
	}
	if payload.IsAnonymous {
		updateRole = GUEST
		updateEmail = pgtype.Text{
			String: "",
			Valid:  false,
		}
	}

	existingUser, err := GetUserById(dbe, idFromPayload)
	// user chưa tồn tại
	if err != nil {
		avtCode, name := genRandomIdentity(adjs, avtCodes)
		err := dbe.q.InsertUser(dbe.c, db.InsertUserParams{
			ID:      idFromPayload,
			Role:    string(updateRole),
			Name:    name,
			AvtCode: string(avtCode),
			Email:   updateEmail,
		})
		if err != nil {
			return User{}, err
		}

		createdUser, err := GetUserById(dbe, idFromPayload)
		if err != nil {
			return User{}, err
		}
		return createdUser, nil
	}

	if existingUser.Role == updateRole && existingUser.Email == updateEmail {
		return existingUser, nil
	}

	err = dbe.q.UpdateUser(dbe.c, db.UpdateUserParams{
		ID:      existingUser.ID,
		Role:    string(updateRole),
		Name:    existingUser.Name,
		AvtCode: string(existingUser.AvtCode),
		Email:   updateEmail,
	})
	if err != nil {
		return User{}, err
	}

	syncedUser, err := dbe.q.GetUserById(dbe.c, idFromPayload)
	if err != nil {
		return User{}, err
	}
	return ToUser(syncedUser), nil
}

var adjs = []string{
	"Lazy",
	"Angry",
	"Silly",
	"Sleepy",
	"Chubby",
	"Happy",
	"Fluffy",
	"Fancy",
	"Lucky",
	"Grumpy",
}

var avtCodes = []UserAvtCode{
	BUNNY,
	KITTEN,
	GRIZZLE,
	HAMSTER,
	MONKEY,
}

func getNameFromAdjAndAvtCode(adj string, avtCode UserAvtCode) string {
	return adj + shared.CapitalizeString(string(avtCode)) + strconv.Itoa(shared.RandomInt(100, 999))
}

func genRandomIdentity(adjs []string, avtCodes []UserAvtCode) (UserAvtCode, string) {
	randAvtCode := avtCodes[shared.RandomInt(0, len(avtCodes)-1)]
	randAdj := adjs[shared.RandomInt(0, len(adjs)-1)]
	randName := getNameFromAdjAndAvtCode(randAdj, randAvtCode)
	return randAvtCode, randName
}
