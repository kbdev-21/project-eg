package domain

import (
	"backend/src/db"
	"backend/src/shared"
	"context"
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

func (a *AppState) GetUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	dbU, err := a.q.GetUserById(ctx, id)
	return ToUser(dbU), err
}

type UpdateUserReq struct {
	Name    string `json:"name" validate:"required,min=2,max=24"`
	AvtCode string `json:"avtCode" validate:"required,oneof=BUNNY KITTEN GRIZZLE HAMSTER MONKEY"` // from UserAvtCode
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
		avtCode, name := genRandomIdentity(adjs, avtCodes)
		err := a.q.InsertUser(ctx, db.InsertUserParams{
			ID:      idFromPayload,
			Role:    string(updateRole),
			Name:    name,
			AvtCode: string(avtCode),
			Email:   updateEmail,
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
		ID:      existingUser.ID,
		Role:    string(updateRole),
		Name:    existingUser.Name,
		AvtCode: string(existingUser.AvtCode),
		Email:   updateEmail,
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
