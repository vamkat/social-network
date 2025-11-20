package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *UserService) RegisterUser(ctx context.Context, req RegisterUserRequest) (UserId, error) {

	// convert date
	dobTime, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return 0, ErrInvalidDateFormat
	}

	dob := pgtype.Date{
		Time:  dobTime,
		Valid: true,
	}

	//hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	var newId int64

	err = s.runTx(ctx, func(q *sqlc.Queries) error {

		// Insert user
		userId, err := q.InsertNewUser(ctx, sqlc.InsertNewUserParams{
			Username:      req.Username,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			DateOfBirth:   dob,
			Avatar:        req.Avatar,
			AboutMe:       req.About,
			ProfilePublic: req.Public,
		})
		if err != nil {
			return err //TODO check how to return correct error
		}
		newId = userId

		// Insert auth
		return q.InsertNewUserAuth(ctx, sqlc.InsertNewUserAuthParams{
			UserID:       userId,
			Email:        req.Email,
			PasswordHash: passwordHash,
		})
	})

	if err != nil {
		return 0, err //TODO check how to return correct error
	}

	return UserId(newId), nil

}

func (s *UserService) LoginUser(ctx context.Context, req LoginReq) (User, error) {
	var u User

	err := s.runTx(ctx, func(q *sqlc.Queries) error {
		row, err := q.GetUserForLogin(ctx, req.Identifier)
		if err != nil {
			return err
		}

		u = User{
			UserId:   row.ID,
			Username: row.Username,
			Avatar:   row.Avatar,
			Public:   row.ProfilePublic,
		}

		if !checkPassword(row.PasswordHash, req.Password) {
			q.IncrementFailedLoginAttempts(ctx, row.ID)
			return err
		}
		q.ResetFailedLoginAttempts(ctx, u.UserId)
		return nil
	})

	if err != nil {
		return User{}, ErrWrongCredentials
	}

	//TODO what happens when eg failed login attempts > 3? Add logic?

	return u, nil
}

func UpdateUserPassword() {
	//called with user_id, old password, new password_hash, salt
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserPassword
}

func UpdateUserEmail() {
	//called with user_id, new email
	//returns success or error
	//request needs to come from same user
	//---------------------------------------------------------------------
	//UpdateUserEmail
}
