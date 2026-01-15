//go:build integration

package user

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/Avito-courses/l11-examples/pkg/db"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	conn   *pgxpool.Pool
	repo   *Repository
	userID int64
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	err := godotenv.Load("../../../.env")
	s.Require().NoError(err)

	s.Require().NoError(err)
	s.conn = db.MustInitDB()
	s.repo = NewUserRepository(s.conn)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.conn.Close()
}

func (s *UserRepositoryTestSuite) SetupTest() {
	ctx := context.Background()

	err := s.conn.QueryRow(
		ctx,
		`INSERT INTO users (name, phone, rating) VALUES ($1, $2, $3) RETURNING id`,
		"John Doe",
		"+79000000000",
		135,
	).Scan(&s.userID)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()

	_, err := s.conn.Exec(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		s.userID,
	)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()

	user, err := s.repo.GetByID(ctx, int(s.userID))

	s.Require().NoError(err)
	s.Equal(s.userID, int64(user.ID))
	s.Equal("John Doe", user.Name)
	s.Equal("+79000000000", user.Phone)
	s.Equal(135, user.Rating)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
