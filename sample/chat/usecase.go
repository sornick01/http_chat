package chat

import (
	"context"
	"sample/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessString string) (*models.User, error)

	GetPrivateMessages(ctx context.Context, user *models.User) ([]*models.Message, error)
	GetGlobalMessages(ctx context.Context) ([]*models.Message, error)
	AddPrivateMessage(ctx context.Context, user *models.User, recipient, textMessage string) error
	AddGlobalMessage(ctx context.Context, user *models.User, textMessage string) error
}
