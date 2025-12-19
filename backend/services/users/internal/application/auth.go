package application

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/services/users/internal/db/sqlc"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgtype"
)

func (app *Application) RegisterUser(ctx context.Context, req models.RegisterUserRequest) (models.User, error) {
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

	queries, commit, rollback, err := app.db.TxQueries(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback(ctx)

	// Insert user
	userId, err := queries.InsertNewUser(ctx, sqlc.InsertNewUserParams{
		Username:      req.Username.String(),
		FirstName:     req.FirstName.String(),
		LastName:      req.LastName.String(),
		DateOfBirth:   dob,
		AvatarID:      req.AvatarId.Int64(),
		AboutMe:       req.About.String(),
		ProfilePublic: req.Public,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert new user: %w", err) //TODO check how to return correct error
	}
	newId = ct.Id(userId)

	// Insert auth
	err = queries.InsertNewUserAuth(ctx, sqlc.InsertNewUserAuthParams{
		UserID:       newId.Int64(),
		Email:        req.Email.String(),
		PasswordHash: req.Password.String(),
	})
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert new user auth: %w", err)
	}

	if err != nil {
		return models.User{}, err //TODO check how to return correct error
	}

	commit(ctx)

	return models.User{
		UserId:   newId,
		Username: req.Username,
		AvatarId: req.AvatarId,
	}, nil

}
func (app *Application) LoginUser(ctx context.Context, req models.LoginRequest) (models.User, error) {
	var u models.User

	if err := ct.ValidateStruct(req); err != nil {
		return u, err
	}

	queries, commit, rollback, err := app.db.TxQueries(ctx)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback(ctx)

	row, err := queries.GetUserForLogin(ctx, sqlc.GetUserForLoginParams{
		Username:     req.Identifier.String(),
		PasswordHash: req.Password.String(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrWrongCredentials
		}
		return models.User{}, err
	}

	u = models.User{
		UserId:   ct.Id(row.ID),
		Username: ct.Username(row.Username),
		AvatarId: ct.Id(row.AvatarID),
	}

	commit(ctx)
	return u, nil
}

func (app *Application) UpdateUserPassword(ctx context.Context, req models.UpdatePasswordRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	queries, commit, rollback, err := app.db.TxQueries(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback(ctx)

	row, err := queries.GetUserPassword(ctx, req.UserId.Int64())
	if err != nil {
		return err
	}

	if !checkPassword(row, req.OldPassword.String()) {
		return ErrNotAuthorized
	}

	err = queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		UserID:       req.UserId.Int64(),
		PasswordHash: req.NewPassword.String(),
	})
	if err != nil {
		return err
	}

	commit(ctx)
	return nil
}

func (app *Application) UpdateUserEmail(ctx context.Context, req models.UpdateEmailRequest) error {
	if err := ct.ValidateStruct(req); err != nil {
		return err
	}

	queries, commit, rollback, err := app.db.TxQueries(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer rollback(ctx)

	err = queries.UpdateUserEmail(ctx, sqlc.UpdateUserEmailParams{
		UserID: req.UserId.Int64(),
		Email:  req.Email.String(),
	})
	if err != nil {
		return err
	}

	commit(ctx)
	return nil
}
