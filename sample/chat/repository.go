package chat

import (
	"context"
	"sample/models"
)

type Repo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username string) (*models.User, error)

	AddGlobalMessage(ctx context.Context, message *models.Message) error
	AddPrivateMessage(ctx context.Context, recipient string, message *models.Message) error
	GetPrivateMessages(ctx context.Context, user *models.User) ([]*models.Message, error)
	GetGlobalMessages(ctx context.Context) ([]*models.Message, error)
}
