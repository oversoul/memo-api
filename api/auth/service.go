package auth

import (
	"context"
	"fmt"
	"strings"

	"memo/pkg/security"
)

type (
	infoRequest struct {
		Name string `json:"name" validate:"required"`
	}

	loginRequest struct {
		Email    string `json:"email" validate:"required|email"`
		Password string `json:"password" validate:"required"`
	}

	registerRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required|email"`
		Password string `json:"password" validate:"required|min:6"`
	}

	registerResponse struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Token string `json:"token"`
		Image string `json:"image"`
	}

	authStore struct {
		repository AuthRepository
	}
)

func (r *registerRequest) AsLogin() *loginRequest {
	return &loginRequest{Email: r.Email, Password: r.Password}
}

type AuthStore interface {
	FindToken(string, context.Context) (string, error)
	CreateUser(*registerRequest, context.Context) error
	DeleteToken(string, context.Context) error
	GetUserById(id string, ctx context.Context) (*AuthUser, error)
	Authenticate(request *loginRequest, ctx context.Context) (*registerResponse, error)
	UpdateUserInfo(u *AuthUser, ctx context.Context) error
}

func NewStore(repository AuthRepository) AuthStore {
	return &authStore{repository}
}

func (s *authStore) GetUserById(id string, ctx context.Context) (*AuthUser, error) {
	return s.repository.FindUserById(id, ctx)
}

func (s *authStore) CreateUser(req *registerRequest, ctx context.Context) error {
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("Error hashing password: %s", err)
	}

	user := AuthUser{Name: req.Name, Email: req.Email, Password: hash}

	if err := s.repository.InsertUser(user, ctx); err != nil {
		return fmt.Errorf("Error creating user: %s", err.Error())
	}

	return nil
}

func (s *authStore) FindToken(token string, ctx context.Context) (string, error) {
	parts := strings.Split(token, "|")

	dbToken, err := s.repository.FindToken(parts[0], ctx)
	if err != nil {
		return "", fmt.Errorf("Error Access Token %s", err.Error())
	}

	if security.HashEquals(dbToken.Token, parts[1]) {
		return dbToken.UserId.Hex(), nil
	}

	return "", fmt.Errorf("Token not valid")
}

func (s *authStore) Authenticate(request *loginRequest, ctx context.Context) (*registerResponse, error) {
	u, err := s.repository.FindUserByEmail(request.Email, ctx)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the saved hash.
	if err := security.CompareHashToPassword(u.Password, request.Password); err != nil {
		return nil, err
	}

	// remove old tokens
	if err := s.repository.DeleteUserTokens(u.ID.Hex(), ctx); err != nil {
		return nil, fmt.Errorf("Couldn't delete tokens: %s", err.Error())
	}

	tok := security.GenerateTokenString()

	// save to db...
	lastId, err := s.repository.InsertToken(u.ID.Hex(), "API TOKEN", tok.Token, ctx)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create token: %s", err.Error())
	}

	// expires_at: time.Now(),
	return &registerResponse{
		Name:  u.Name,
		Email: u.Email,
		Image: u.Image,
		Token: fmt.Sprintf("%s|%s", lastId, tok.Plain),
	}, nil
}

func (s *authStore) DeleteToken(token string, ctx context.Context) error {
	return s.repository.DeleteToken(token, ctx)
}

func (s *authStore) UpdateUserInfo(u *AuthUser, ctx context.Context) error {
	return s.repository.UpdateUserInfo(u, ctx)
}
