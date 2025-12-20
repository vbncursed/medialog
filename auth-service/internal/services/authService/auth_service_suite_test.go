package authService_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService/mocks"
)

type AuthServiceSuite struct {
	suite.Suite

	ctx context.Context
	st  *mocks.Storage
	svc *authService.AuthService
}

func (s *AuthServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.st = mocks.NewStorage(s.T())
	s.svc = newTestService(s.st)
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
