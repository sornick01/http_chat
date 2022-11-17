package localcache

import (
	"context"
	"errors"
	"sample/chat"
	"sample/models"
)

type LocalStorage struct {
	userStorage           *UserStorage
	globalMessageStorage  *GlobalMessageStorage
	privateMessageStorage *PrivateMessageStorage
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		userStorage:           NewUserStorage(),
		globalMessageStorage:  NewGlobalMessageStorage(),
		privateMessageStorage: NewPrivateMessageStorage(),
	}
}

func (l *LocalStorage) CreateUser(ctx context.Context, user *models.User) error {
	l.userStorage.mutex.Lock()
	defer l.userStorage.mutex.Unlock()
	l.userStorage.idCounter++
	user.ID = l.userStorage.idCounter
	l.userStorage.users[user.ID] = user

	return nil
}

func (l *LocalStorage) GetUser(ctx context.Context, username string) (*models.User, error) {
	l.userStorage.mutex.Lock()
	defer l.userStorage.mutex.Unlock()

	for _, user := range l.userStorage.users {
		if username == user.Username {
			return user, nil
		}
	}

	return nil, chat.ErrUserNotFound
}

func (l *LocalStorage) AddGlobalMessage(ctx context.Context, message *models.Message) error {
	l.globalMessageStorage.mutex.Lock()
	defer l.globalMessageStorage.mutex.Unlock()
	l.globalMessageStorage.messages = append(l.globalMessageStorage.messages, message)

	return nil
}

func (l *LocalStorage) AddPrivateMessage(ctx context.Context, recipient string, message *models.Message) error {
	l.privateMessageStorage.mutex.Lock()
	defer l.privateMessageStorage.mutex.Unlock()
	for k, _ := range l.privateMessageStorage.messages {
		if k == recipient {
			l.privateMessageStorage.messages[recipient] = append(l.privateMessageStorage.messages[recipient], message)
			return nil
		}
	}

	return errors.New("no such user")
}

func (l *LocalStorage) GetGlobalMessages(ctx context.Context) ([]*models.Message, error) {
	l.globalMessageStorage.mutex.Lock()
	defer l.globalMessageStorage.mutex.Unlock()

	return l.globalMessageStorage.messages, nil
}

func (l *LocalStorage) GetPrivateMessages(ctx context.Context, user *models.User) ([]*models.Message, error) {
	l.privateMessageStorage.mutex.Lock()
	defer l.privateMessageStorage.mutex.Unlock()
	return l.privateMessageStorage.messages[user.Username], nil
}
