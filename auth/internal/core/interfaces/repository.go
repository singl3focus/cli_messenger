package interfaces

import (
	"context"

	"github.com/singl3focus/cli_messenger/auth/internal/core/models"
)

type Repostory interface {
	UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, id int64, name, email *string) error
	DeleteUser(ctx context.Context, id int64) error
}
