package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"rest-jwt/internal/app/models"
	"rest-jwt/internal/app/repository/postgresql"
)

type Service interface {
	GenerateToken(userID, clientIP string) (string, string, error)
	RefreshToken(userID, refreshToken, clientIP string) (string, string, error)
}

type AuthService struct {
	repo   postgresql.Repository
	jwtKey string
}

func New(repo postgresql.Repository, jwtKey string) *AuthService {
	return &AuthService{repo: repo, jwtKey: jwtKey}
}

func (as *AuthService) GenerateToken(userID, clientIP string) (string, string, error) {
	accessToken, err := as.createAccessToken(userID, clientIP)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return "", "", err
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	err = as.repo.SaveRefreshToken(userID, string(hashedRefreshToken), clientIP)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (as *AuthService) RefreshToken(accessToken, refreshToken, clientIP string) (string, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return as.jwtKey, nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid access token")
	}

	if claims.ClientIP != clientIP {
		sendEmailWarning(claims.UserID)
	}

	dbHash, err := as.repo.GetRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(refreshToken))
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	return as.GenerateToken(claims.UserID, clientIP)
}

func (as *AuthService) createAccessToken(userID, clientIP string) (string, error) {
	claims := &models.TokenClaims{
		UserID:   userID,
		ClientIP: clientIP,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(as.jwtKey)
}

func createRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func sendEmailWarning(userID string) {
	fmt.Printf("Warning: IP address changed for user %s\n", userID)
}
