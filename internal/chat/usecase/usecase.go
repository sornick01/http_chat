package usecase

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	chat2 "github.com/sornick01/http_chat/internal/chat"
	models2 "github.com/sornick01/http_chat/internal/models"
	"time"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	User *models2.User
}

type Chat struct {
	repo           chat2.Repo
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewChat(
	repo chat2.Repo,
	hashSalt string,
	signingKey []byte,
	expireDuration time.Duration) *Chat {
	return &Chat{
		repo:           repo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration}
}

func (c *Chat) SignUp(ctx context.Context, username, password string) error {
	_, err := c.repo.GetUser(ctx, username)
	if err == nil {
		return errors.New("username occupied")
	}
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(c.hashSalt))

	user := &models2.User{
		Username: username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return c.repo.CreateUser(ctx, user)
}

func (c *Chat) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(c.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := c.repo.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	if user.Password != password {
		return "", chat2.ErrUserNotFound
	}

	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(c.signingKey)
}

func (c *Chat) ParseToken(ctx context.Context, accessString string) (*models2.User, error) {
	token, err := jwt.ParseWithClaims(accessString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return c.signingKey, nil
	})

	if err != nil {
		return nil, chat2.ErrInvalidAccessToken
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, chat2.ErrInvalidAccessToken
}

func (c *Chat) GetPrivateMessages(ctx context.Context, user *models2.User) ([]*models2.Message, error) {
	return c.repo.GetPrivateMessages(ctx, user)
}

func (c *Chat) GetGlobalMessages(ctx context.Context) ([]*models2.Message, error) {
	return c.repo.GetGlobalMessages(ctx)
}

func (c *Chat) AddPrivateMessage(ctx context.Context, user *models2.User, recipient, textMessage string) error {
	message := &models2.Message{
		Author: user.Username,
		Text:   textMessage,
	}

	return c.repo.AddPrivateMessage(ctx, recipient, message)
}

func (c *Chat) AddGlobalMessage(ctx context.Context, user *models2.User, textMessage string) error {
	message := &models2.Message{
		Author: user.Username,
		Text:   textMessage,
	}

	return c.repo.AddGlobalMessage(ctx, message)
}
