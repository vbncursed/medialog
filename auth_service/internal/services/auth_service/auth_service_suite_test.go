package auth_service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service/mocks"
)

type AuthServiceSuite struct {
	suite.Suite

	ctx context.Context
	st  *mocks.AuthStorage
	svc *auth_service.AuthService
}

func (s *AuthServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.st = mocks.NewAuthStorage(s.T())
	s.svc = newTestService(s.st)
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
