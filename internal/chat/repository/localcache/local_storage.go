package localcache

import (
	"context"
	"github.com/sornick01/http_chat/internal/chat"
	models2 "github.com/sornick01/http_chat/internal/models"
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

func (l *LocalStorage) CreateUser(ctx context.Context, user *models2.User) error {
	l.userStorage.mutex.Lock()
	defer l.userStorage.mutex.Unlock()
	l.userStorage.idCounter++
	user.ID = l.userStorage.idCounter
	l.userStorage.users[user.ID] = user

	return nil
}

func (l *LocalStorage) GetUser(ctx context.Context, username string) (*models2.User, error) {
	l.userStorage.mutex.Lock()
	defer l.userStorage.mutex.Unlock()

	for _, user := range l.userStorage.users {
		if username == user.Username {
			return user, nil
		}
	}

	return nil, chat.ErrUserNotFound
}

func (l *LocalStorage) AddGlobalMessage(ctx context.Context, message *models2.Message) error {
	l.globalMessageStorage.mutex.Lock()
	defer l.globalMessageStorage.mutex.Unlock()
	l.globalMessageStorage.messages = append(l.globalMessageStorage.messages, message)

	return nil
}

func (l *LocalStorage) AddPrivateMessage(ctx context.Context, recipient string, message *models2.Message) error {
	l.privateMessageStorage.mutex.Lock()
	defer l.privateMessageStorage.mutex.Unlock()
	if _, inMap := l.privateMessageStorage.messages[recipient]; inMap {
		l.privateMessageStorage.messages[recipient] = append(l.privateMessageStorage.messages[recipient], message)
	} else {
		l.privateMessageStorage.messages[recipient] = []*models2.Message{message}
	}
	//for k, _ := range l.privateMessageStorage.messages {
	//	if k == recipient {
	//		l.privateMessageStorage.messages[recipient] = append(l.privateMessageStorage.messages[recipient], message)
	//		return nil
	//	}
	//}

	return nil
}

func (l *LocalStorage) GetGlobalMessages(ctx context.Context) ([]*models2.Message, error) {
	l.globalMessageStorage.mutex.Lock()
	defer l.globalMessageStorage.mutex.Unlock()

	return l.globalMessageStorage.messages, nil
}

func (l *LocalStorage) GetPrivateMessages(ctx context.Context, user *models2.User) ([]*models2.Message, error) {
	l.privateMessageStorage.mutex.Lock()
	defer l.privateMessageStorage.mutex.Unlock()
	return l.privateMessageStorage.messages[user.Username], nil
}
