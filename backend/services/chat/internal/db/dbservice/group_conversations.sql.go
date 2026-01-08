package dbservice

import (
	"context"
	"errors"
	"fmt"
	ce "social-network/shared/go/commonerrors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"

	"github.com/jackc/pgx/v5"
)

// Initiates the conversation by groupId. Group Id should be a not null value.
// If conversation already exists then the conversation id is returned.
// Use as a preparation for adding members on conversation or adding messages in
// in group conversations.
func (q *Queries) CreateGroupConv(ctx context.Context,
	groupId ct.Id) (convId ct.Id, err error) {
	row := q.db.QueryRow(ctx, createGroupConv, groupId)
	err = row.Scan(&convId)
	if err != nil {
		return 0, ce.New(ce.ErrInternal,
			fmt.Errorf("failed to create group conversation: %w", err),
			fmt.Sprintf("conversation id: %d", convId),
		)
	}
	return convId, err
}

func (q *Queries) CreateNewGroupMessage(ctx context.Context,
	arg md.CreateGroupMsgReq) (msg md.PrivateMsg, err error) {
	input := fmt.Sprintf("arg: %#v", arg)

	row := q.db.QueryRow(ctx,
		createGroupMessage,
		arg.GroupId,
		arg.SenderId,
		arg.MessageText,
	)

	err = row.Scan(
		&msg.Id,
		&msg.ConversationID,
		&msg.Sender.UserId,
		&msg.MessageText,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&msg.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return msg, ce.New(ce.ErrInvalidArgument, err, input)
		}
		return msg, ce.New(ce.ErrInternal, err, input)
	}
	return msg, err
}

// TODO: GetBefore and after messages
