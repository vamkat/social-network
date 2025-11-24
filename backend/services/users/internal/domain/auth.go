package userservice

import (
	"context"
	"social-network/services/users/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO make repo layer to handle sqlc to domain model conversions

func (s *UserService) RegisterUser(ctx context.Context, req RegisterUserRequest) (User, error) {
	//if no username assign full name
	if req.Username == "" {
		req.Username = req.FirstName + "_" + req.LastName
	}

	// convert date
	dob := pgtype.Date{
		Time:  req.DateOfBirth,
		Valid: true,
	}

	//hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return User{}, err
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
			UserID:       newId,
			Email:        req.Email,
			PasswordHash: passwordHash,
		})
	})

	if err != nil {
		return User{}, err //TODO check how to return correct error
	}

	return User{
		UserId:   newId,
		Username: req.Username,
		Avatar:   req.Avatar,
	}, nil

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
		}

		if !checkPassword(row.PasswordHash, req.Password) {
			return err
		}
		return nil
	})

	if err != nil {
		return User{}, ErrWrongCredentials
	}

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

	err := s.db.UpdateUserEmail(ctx, sqlc.UpdateUserEmailParams{
		UserID: req.UserId,
		Email:  req.Email,
	})
	if err != nil {
		return err
	}
	return nil
}
