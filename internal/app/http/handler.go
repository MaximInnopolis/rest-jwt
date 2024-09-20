package http

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"rest-jwt/internal/app/api"
)

type Handler struct {
	service api.Service
}

func New(service api.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.UserID == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)

	accessToken, refreshToken, err := h.service.GenerateToken(req.UserID, clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.AccessToken == "" || req.RefreshToken == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)

	accessToken, refreshToken, err := h.service.RefreshToken(req.AccessToken, req.RefreshToken, clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func getClientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/auth/token", h.GenerateToken).Methods("POST")
	r.HandleFunc("/auth/refresh", h.RefreshToken).Methods("POST")
}
