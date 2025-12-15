package application

import (
	"context"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Application) RegisterUser(ctx context.Context, req models.RegisterUserRequest) (models.User, error) {
	if err := ct.ValidateStruct(req); err != nil {
		return models.User{}, err
	}

	//if no username assign full name
	if req.Username == "" {
		req.Username = ct.Username(string(req.FirstName) + "_" + string(req.LastName))
	}

	// convert date
	dob := pgtype.Date{
		Time:  req.DateOfBirth.Time(),
		Valid: true,
	}

	var newId ct.Id

	err := s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {

		// Insert user
		userId, err := q.InsertNewUser(ctx, sqlc.InsertNewUserParams{
			Username:      req.Username.String(),
			FirstName:     req.FirstName.String(),
			LastName:      req.LastName.String(),
			DateOfBirth:   dob,
			AvatarID:      req.AvatarId.Int64(),
			AboutMe:       req.About.String(),
			ProfilePublic: req.Public,
		})
		if err != nil {
			return err //TODO check how to return correct error
		}
		newId = ct.Id(userId)

		// Insert auth
		return q.InsertNewUserAuth(ctx, sqlc.InsertNewUserAuthParams{
			UserID:       newId.Int64(),
			Email:        req.Email.String(),
			PasswordHash: req.Password.String(),
		})
	})

	if err != nil {
		return models.User{}, err //TODO check how to return correct error
	}

	return models.User{
		UserId:   newId,
		Username: req.Username,
		AvatarId: req.AvatarId,
	}, nil

}

func (s *Application) LoginUser(ctx context.Context, req models.LoginRequest) (models.User, error) {
	var u models.User

	if err := ct.ValidateStruct(req); err != nil {
		return u, err
	}

	err := s.txRunner.RunTx(ctx, func(q sqlc.Querier) error {
		row, err := q.GetUserForLogin(ctx, req.Identifier.String())
		if err != nil {
			return err
		}

		u = models.User{
			UserId:   ct.Id(row.ID),
			Username: ct.Username(row.Username),
			AvatarId: ct.Id(row.AvatarID),
		}

		if !checkPassword(row.PasswordHash, req.Password.String()) {
			return ErrWrongCredentials
		}
		return nil
	})

	if err != nil {
		return models.User{}, ErrWrongCredentials
	}

	return u, nil
}

func (s *Application) UpdateUserPassword(ctx context.Context, req models.UpdatePasswordRequest) error {

	//TODO think whether transaction is needed here
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	if !checkPassword(req.NewPassword.String(), req.OldPassword.String()) {
		return ErrNotAuthorized
	}

	err := s.db.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		UserID:       req.UserId.Int64(),
		PasswordHash: req.NewPassword.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Application) UpdateUserEmail(ctx context.Context, req models.UpdateEmailRequest) error {

	//TODO validate email
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	err := s.db.UpdateUserEmail(ctx, sqlc.UpdateUserEmailParams{
		UserID: req.UserId.Int64(),
		Email:  req.Email.String(),
	})
	if err != nil {
		return err
	}
	return nil
}
