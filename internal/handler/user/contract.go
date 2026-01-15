//go:generate mockgen -source=contract.go -destination=./mocks/repo_mock.go -package=mocks
package user

import (
	"context"

	"github.com/Avito-courses/l11-examples/internal/model/user"
)

type userRepo interface {
	GetByID(ctx context.Context, id int) (*user.User, error)
}
