package domain

import (
	"backend/src/db"
	"backend/src/shared"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
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
	KITTY   UserAvtCode = "KITTY"
	TEDDY   UserAvtCode = "TEDDY"
	HAMSTER UserAvtCode = "HAMSTER"
	MONKEY  UserAvtCode = "MONKEY"
	PIGGY   UserAvtCode = "PIGGY"
)

func ToUser(u db.User) User {
	return User{
		User:    u,
		Role:    UserRole(u.Role),
		AvtCode: UserAvtCode(u.AvtCode),
	}
}

func (a *AppState) GetUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	dbU, err := a.q.GetUserById(ctx, id)
	return ToUser(dbU), err
}

func (a *AppState) GetUserByName(ctx context.Context, name string) (User, error) {
	dbU, err := a.q.GetUserByName(ctx, name)
	return ToUser(dbU), err
}

type UpdateUserReq struct {
	Name    string `json:"name" validate:"required,min=2,max=24"`
	AvtCode string `json:"avtCode" validate:"required,oneof=BUNNY KITTY TEDDY HAMSTER MONKEY PIGGY"` // from UserAvtCode
}

func (a *AppState) UpdateUserInfo(ctx context.Context, id pgtype.UUID, req UpdateUserReq) (User, error) {
	existingUser, err := a.GetUserById(ctx, id)
	if err != nil {
		return User{}, err
	}

	err = a.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:      id,
		Role:    string(existingUser.Role),
		Name:    req.Name,
		NormalizedName: shared.NormalizedText(req.Name),
		AvtCode: req.AvtCode,
		Email:   existingUser.Email,
	})
	if err != nil {
		return User{}, err
	}

	updatedUser, err := a.GetUserById(ctx, id)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func (a *AppState) SyncUserFromTokenPayload(ctx context.Context, payload TokenPayload) (User, error) {
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

	existingUser, err := a.GetUserById(ctx, idFromPayload)
	// user chưa tồn tại
	if err != nil {
		avtCode, name := a.GenRandomUniqueIdentity(ctx, adjs, avtCodes)
		err := a.q.InsertUser(ctx, db.InsertUserParams{
			ID:             idFromPayload,
			Role:           string(updateRole),
			Name:           name,
			NormalizedName: shared.NormalizedText(name),
			AvtCode:        string(avtCode),
			Email:          updateEmail,
		})
		if err != nil {
			return User{}, err
		}

		createdUser, err := a.GetUserById(ctx, idFromPayload)
		if err != nil {
			return User{}, err
		}
		return createdUser, nil
	}

	if existingUser.Role == updateRole && existingUser.Email == updateEmail {
		return existingUser, nil
	}

	err = a.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:             existingUser.ID,
		Role:           string(updateRole),
		Name:           existingUser.Name,
		NormalizedName: existingUser.NormalizedName,
		AvtCode:        string(existingUser.AvtCode),
		Email:          updateEmail,
	})
	if err != nil {
		return User{}, err
	}

	syncedUser, err := a.q.GetUserById(ctx, idFromPayload)
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
	"Fancy",
	"Lucky",
	"Greedy",
	"Funny",
	"Sneaky",
	"Tiny",
	"Cheeky",
	"Naughty",
	"Noisy",
}

var avtCodes = []UserAvtCode{
	BUNNY,
	KITTY,
	TEDDY,
	HAMSTER,
	MONKEY,
	PIGGY,
}

func (a *AppState) GenRandomUniqueIdentity(ctx context.Context, adjs []string, avtCodes []UserAvtCode) (UserAvtCode, string) {
	const MAX_ATTEMPTS = 10

	var randName string
	var randAvtCode UserAvtCode
	isSuccess := false

	for i := range MAX_ATTEMPTS {
		start, end := 100, 999
		if i >= MAX_ATTEMPTS/2 {
			start, end = 1000, 9999
		}

		randAvtCode = avtCodes[shared.RandomInt(0, len(avtCodes)-1)]
		randAdj := adjs[shared.RandomInt(0, len(adjs)-1)]
		randDigit := shared.RandomInt(start, end)
		randName = getNameFromAdjAndAvtCode(randAdj, randAvtCode, randDigit)

		_, err := a.GetUserByName(ctx, randName)
		if err != nil {
			isSuccess = true
			break
		}
	}

	if !isSuccess {
		randAvtCode = avtCodes[shared.RandomInt(0, len(avtCodes)-1)]
		randAdj := adjs[shared.RandomInt(0, len(adjs)-1)]
		randDigit := shared.RandomInt(1000_0000, 9999_9999)
		randName = getNameFromAdjAndAvtCode(randAdj, randAvtCode, randDigit)
	}

	return randAvtCode, randName
}

func getNameFromAdjAndAvtCode(adj string, avtCode UserAvtCode, digit int) string {
	return adj + shared.CapitalizeString(string(avtCode)) + strconv.Itoa(digit)
}
