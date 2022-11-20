package localcache

import (
	"github.com/sornick01/http_chat/internal/models"
	"sync"
)

type UserStorage struct {
	idCounter int
	users     map[int]*models.User
	mutex     *sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[int]*models.User),
		mutex: new(sync.Mutex),
	}
}
