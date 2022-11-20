package chat

import (
	"context"
	models2 "github.com/sornick01/http_chat/internal/models"
)

type Repo interface {
	CreateUser(ctx context.Context, user *models2.User) error
	GetUser(ctx context.Context, username string) (*models2.User, error)

	AddGlobalMessage(ctx context.Context, message *models2.Message) error
	AddPrivateMessage(ctx context.Context, recipient string, message *models2.Message) error
	GetPrivateMessages(ctx context.Context, user *models2.User) ([]*models2.Message, error)
	GetGlobalMessages(ctx context.Context) ([]*models2.Message, error)
}
