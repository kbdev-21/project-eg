package auth

import (
	"backend/src/db"
	"backend/src/shared"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

func SyncUserFromTokenPayload(c fiber.Ctx, q *db.Queries, p TokenPayload) (db.User, error) {
	idFromPayload := shared.ForceParseStringToUuid(p.Sub)
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
		avtCode, name := genRandomIdentity()
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
