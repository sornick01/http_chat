package localcache

import (
	"sample/models"
	"sync"
)

type GlobalMessageStorage struct {
	messages []*models.Message
	mutex    *sync.Mutex
}

func NewGlobalMessageStorage() *GlobalMessageStorage {
	return &GlobalMessageStorage{
		messages: make([]*models.Message, 0),
		mutex:    new(sync.Mutex),
	}
}

type PrivateMessageStorage struct {
	messages map[string][]*models.Message
	mutex    *sync.Mutex
}

func NewPrivateMessageStorage() *PrivateMessageStorage {
	return &PrivateMessageStorage{
		messages: make(map[string][]*models.Message),
		mutex:    new(sync.Mutex),
	}
}
