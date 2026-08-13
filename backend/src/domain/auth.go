package domain

import (
	"backend/src/db"
	"backend/src/shared"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func SyncUserFromTokenPayload(c context.Context, q *db.Queries, p TokenPayload) (db.User, error) {
	idFromPayload, err := shared.ParseStringToUuid(p.Sub)
	if err != nil {
		return db.User{}, err
	}

	updateRole := db.UserRoleUSER
	updateEmail := pgtype.Text{
		String: p.Email,
		Valid:  true,
	}
	if p.IsAnonymous {
		updateRole = db.UserRoleGUEST
		updateEmail = pgtype.Text{
			String: "",
			Valid:  false,
		}
	}

	existingUser, err := q.GetUserById(c, idFromPayload)
	// user chưa tồn tại
	if err != nil {
		avtCode, name := genRandomIdentity(adjs, avtCodes)
		err := q.InsertUser(c, db.InsertUserParams{
			ID:      idFromPayload,
			Role:    updateRole,
			Name:    name,
			AvtCode: avtCode,
			Email:   updateEmail,
		})
		if err != nil {
			return db.User{}, err
		}

		createdUser, err := q.GetUserById(c, idFromPayload)
		if err != nil {
			return db.User{}, err
		}
		return createdUser, nil
	}

	if existingUser.Role == updateRole && existingUser.Email == updateEmail {
		return existingUser, nil
	}

	err = q.UpdateUser(c, db.UpdateUserParams{
		ID:      existingUser.ID,
		Role:    updateRole,
		Name:    existingUser.Name,
		AvtCode: existingUser.AvtCode,
		Email:   updateEmail,
	})
	if err != nil {
		return db.User{}, err
	}

	syncedUser, err := q.GetUserById(c, idFromPayload)
	if err != nil {
		return db.User{}, err
	}
	return syncedUser, nil
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

var avtCodes = []db.UserAvtCode{
	db.UserAvtCodeBUNNY,
	db.UserAvtCodeKITTEN,
	db.UserAvtCodeGRIZZLE,
	db.UserAvtCodeHAMSTER,
	db.UserAvtCodeMONKEY,
}

func getNameFromAdjAndAvtCode(adj string, avtCode db.UserAvtCode) string {
	return adj + shared.CapitalizeString(string(avtCode))
}

func genRandomIdentity(adjs []string, avtCodes []db.UserAvtCode) (db.UserAvtCode, string) {
	randAvtCode := avtCodes[shared.RandomInt(0, len(avtCodes)-1)]
	randAdj := adjs[shared.RandomInt(0, len(adjs)-1)]
	randName := getNameFromAdjAndAvtCode(randAdj, randAvtCode)
	return randAvtCode, randName
}