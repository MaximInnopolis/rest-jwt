package api

import (
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"rest-jwt/internal/app/models"

	mockrepo "rest-jwt/internal/app/repository/postgresql/mocks"
)

func TestAuthService_GenerateToken(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtKey := "gJp5K#L$m8n3H@jR&tT*qX7pZ!bD%wFy"
	mockRepo := mockrepo.NewMockRepository(ctrl)
	authService := New(mockRepo, jwtKey)

	userID := "testUserID"
	clientIP := "127.0.0.1"

	t.Run("successful token generation", func(t *testing.T) {
		t.Parallel()

		mockRepo.EXPECT().SaveRefreshToken(userID, gomock.Any(), clientIP).Return(nil)

		accessToken, refreshToken, err := authService.GenerateToken(userID, clientIP)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
	})

	t.Run("repository error on saving refresh token", func(t *testing.T) {
		t.Parallel()

		mockRepo.EXPECT().SaveRefreshToken(userID, gomock.Any(), clientIP).Return(errors.New("database error"))

		_, _, err := authService.GenerateToken(userID, clientIP)
		assert.EqualError(t, err, "database error")
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Parallel()
	jwtKey := "gJp5K#L$m8n3H@jR&tT*qX7pZ!bD%wFy"
	userID := "testUserID"
	clientIP := "127.0.0.1"
	refreshToken := "testRefreshToken"

	accessToken, err := createTestAccessToken(userID, clientIP, jwtKey)
	if err != nil {
		t.Fatalf("failed to create test access token: %v", err)
	}

	hashedRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	t.Run("successful token refresh", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepository(ctrl)
		authService := New(mockRepo, jwtKey)

		mockRepo.EXPECT().GetRefreshToken(userID).Return(string(hashedRefreshToken), nil)
		mockRepo.EXPECT().SaveRefreshToken(userID, gomock.Any(), clientIP).Return(nil)

		newAccessToken, newRefreshToken, err := authService.RefreshToken(accessToken, refreshToken, clientIP)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
	})

	t.Run("access token IP mismatch", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockrepo.NewMockRepository(ctrl)
		authService := New(mockRepo, jwtKey)

		mockRepo.EXPECT().GetRefreshToken(userID).Return(string(hashedRefreshToken), nil)
		mockRepo.EXPECT().SaveRefreshToken(userID, gomock.Any(), clientIP).Return(nil)

		differentIP := "192.168.1.1"
		accessTokenWithDifferentIP, err := createTestAccessToken(userID, differentIP, jwtKey)
		assert.NoError(t, err)

		newAccessToken, newRefreshToken, err := authService.RefreshToken(accessTokenWithDifferentIP, refreshToken, clientIP)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
	})
}

func createTestAccessToken(userID, clientIP, jwtKey string) (string, error) {
	claims := &models.TokenClaims{
		UserID:   userID,
		ClientIP: clientIP,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(jwtKey))
}
