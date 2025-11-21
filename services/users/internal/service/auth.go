package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// username becomes FirstName_Lastname
// uuid to send to front
// get time.Time instead of string
// register returns basicUserInfo (id, username, avatar) --public (remove)
// API GATEWAY checks password at least 8 characters,check email
// make token
//add owner to group

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

func (s *UserService) LoginUser(ctx context.Context, req LoginRequest) (User, error) {
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

	//TODO what happens when eg failed login attempts > 3? Add logic? //remove

	return u, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, req UpdatePasswordRequest) error {
	//TODO think whether transaction is needed here

	hashedPassword, err := s.db.GetUserPassword(ctx, req.UserId)
	if err != nil {
		return err
	}

	if !checkPassword(hashedPassword, req.OldPassword) {
		return ErrNotAuthorized
	}

	newPasswordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.db.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		UserID:       req.UserId,
		PasswordHash: newPasswordHash,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserEmail(ctx context.Context, req UpdateEmailRequest) error {
	//reminder: userId always from token
	err := s.db.UpdateUserEmail(ctx, sqlc.UpdateUserEmailParams{
		UserID: req.UserId,
		Email:  req.Email,
	})
	if err != nil {
		return err
	}
	return nil
}
