package grpc

import (
	"context"

	"auth-service/internal/auth/app"
	"auth-service/internal/auth/delivery/grpc/authpb"
	"auth-service/internal/auth/domain"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	service *app.AuthService
}

// NewAuthHandler membuat handler baru
func NewAuthHandler(service *app.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// RegisterUser menerima request gRPC dan mendaftarkan user
func (h *AuthHandler) RegisterUser(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := domain.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	createdUser, err := h.service.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		Id:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

// LoginUser melakukan autentikasi dan menghasilkan token
func (h *AuthHandler) LoginUser(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, err := h.service.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Token: token,
	}, nil
}

// ValidateToken memvalidasi token JWT dan mengembalikan data user
func (h *AuthHandler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	user, err := h.service.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &authpb.ValidateTokenResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
