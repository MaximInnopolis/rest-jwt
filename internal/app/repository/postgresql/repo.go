package postgresql

import (
	"context"

	"rest-jwt/internal/app/repository/database"
)

type Repository interface {
	SaveRefreshToken(userID, refreshToken, clientIP string) error
	GetRefreshToken(userID string) (string, error)
}

type Repo struct {
	db database.Database
}

func New(db database.Database) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveRefreshToken(userID, refreshToken, clientIP string) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, client_ip) VALUES ($1, $2, $3)
              ON CONFLICT (user_id) DO UPDATE SET token_hash = $2, client_ip = $3`
	ctx := context.Background()

	_, err := r.db.GetPool().Exec(ctx, query, userID, refreshToken, clientIP)
	return err
}

func (r *Repo) GetRefreshToken(userID string) (string, error) {
	query := `SELECT token_hash FROM refresh_tokens WHERE user_id = $1`
	ctx := context.Background()
	var token string

	err := r.db.GetPool().QueryRow(ctx, query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
