package chat

import (
	"context"
	models2 "github.com/sornick01/http_chat/internal/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessString string) (*models2.User, error)

	GetPrivateMessages(ctx context.Context, user *models2.User) ([]*models2.Message, error)
	GetGlobalMessages(ctx context.Context) ([]*models2.Message, error)
	AddPrivateMessage(ctx context.Context, user *models2.User, recipient, textMessage string) error
	AddGlobalMessage(ctx context.Context, user *models2.User, textMessage string) error
}
