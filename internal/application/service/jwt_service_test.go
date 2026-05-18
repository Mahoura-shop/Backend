package service

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type jwtKeyManagerStub struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func (k *jwtKeyManagerStub) LoadKeys(_, _ string) error  { return nil }
func (k *jwtKeyManagerStub) GetPrivateKey() *rsa.PrivateKey { return k.private }
func (k *jwtKeyManagerStub) GetPublicKey() *rsa.PublicKey   { return k.public }

func newTestKeyManager(t *testing.T) *jwtKeyManagerStub {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	return &jwtKeyManagerStub{private: key, public: &key.PublicKey}
}

type JWTServiceTestSuite struct {
	suite.Suite
	service    *JWTService
	keyManager *jwtKeyManagerStub
}

func (s *JWTServiceTestSuite) SetupTest() {
	s.keyManager = newTestKeyManager(s.T())
	keysPath := &bootstrap.JWTKeysPath{}
	s.service = NewJWTService(s.keyManager, keysPath)
}

func (s *JWTServiceTestSuite) TestGenerateToken_ReturnsNonEmptyStrings() {
	access, refresh, err := s.service.GenerateToken(42)

	s.NoError(err)
	s.NotEmpty(access)
	s.NotEmpty(refresh)
}

func (s *JWTServiceTestSuite) TestGenerateToken_DifferentTokens() {
	access, refresh, err := s.service.GenerateToken(1)

	s.NoError(err)
	s.NotEqual(access, refresh)
}

func (s *JWTServiceTestSuite) TestValidateToken_Success() {
	access, _, err := s.service.GenerateToken(7)
	s.Require().NoError(err)

	claims, err := s.service.ValidateToken(access)

	s.NoError(err)
	s.NotNil(claims)
	s.Equal(float64(7), claims["sub"])
}

func (s *JWTServiceTestSuite) TestValidateToken_Tampered() {
	access, _, err := s.service.GenerateToken(7)
	s.Require().NoError(err)

	tampered := access + "x"
	_, err = s.service.ValidateToken(tampered)

	s.Error(err)
}

func (s *JWTServiceTestSuite) TestValidateToken_WrongAlgorithm() {
	hmacToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": float64(1),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := hmacToken.SignedString([]byte("secret"))
	s.Require().NoError(err)

	_, err = s.service.ValidateToken(signed)

	s.Error(err)
}

func (s *JWTServiceTestSuite) TestValidateToken_Expired() {
	claims := jwt.MapClaims{
		"sub": float64(1),
		"exp": time.Now().Add(-time.Hour).Unix(),
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.keyManager.GetPrivateKey())
	s.Require().NoError(err)

	_, err = s.service.ValidateToken(signed)

	s.Error(err)
}

func (s *JWTServiceTestSuite) TestGenerateToken_ClaimsContainUserID() {
	access, _, err := s.service.GenerateToken(99)
	s.Require().NoError(err)

	claims, err := s.service.ValidateToken(access)
	s.Require().NoError(err)

	s.Equal(float64(99), claims["sub"])
}

func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}
