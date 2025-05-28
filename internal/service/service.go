package service

import (
	"context"
	"log/slog"

	"github.com/OxytocinGroup/theca-v3/internal/config"
	"github.com/OxytocinGroup/theca-v3/internal/model"
	"github.com/OxytocinGroup/theca-v3/internal/repository"
	jwtauth "github.com/OxytocinGroup/theca-v3/internal/utils/jwt"
	"github.com/OxytocinGroup/theca-v3/internal/vars"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, email, username, password string) error
	Login(ctx context.Context, username, password string) (string, string, error)
	LogoutFromAllSessions(ctx context.Context, userID uint) error
}

type service struct {
	repo repository.Repository
	log  *slog.Logger
	cfg  *config.Config
}

func NewService(repo repository.Repository, log *slog.Logger, cfg *config.Config) Service {
	return &service{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *service) Register(ctx context.Context, email, username, password string) error {
	const op = "service.Register"
	log := s.log.With("op", op)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return err
	}

	user := model.User{
		Email:    email,
		Username: username,
		PassHash: string(hashPassword),
	}

	err = s.repo.Register(ctx, &user)
	if err != nil {
		return err
	}
	log.Debug("user registered", "user", user.ID)
	return nil
}

func (s *service) Login(ctx context.Context, username, password string) (string, string, error) {
	const op = "service.Login"
	log := s.log.With("op", op)

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", vars.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		log.Error("failed to compare password", "error", err)
		return "", "", vars.ErrInvalidPassword
	}

	accessToken, err := jwtauth.GenerateAccessToken(user.ID, s.cfg.JWTAccessSecret)
	if err != nil {
		log.Error("failed to generate access token", "error", err)
		return "", "", err
	}

	refreshToken, err := jwtauth.GenerateRefreshToken(user.ID, user.RefreshTokenVersion, s.cfg.JWTRefreshSecret)
	if err != nil {
		log.Error("failed to generate refresh token", "error", err)
		return "", "", err
	}

	log.Debug("login successful", "user", user.ID)
	return accessToken, refreshToken, nil
}

func (s *service) LogoutFromAllSessions(ctx context.Context, userID uint) error {
	const op = "service.LogoutFromAllSessions"
	log := s.log.With("op", op)

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.RefreshTokenVersion += 1

	err = s.repo.SaveUser(ctx, user)
	if err != nil {
		return err
	}

	log.Debug("logout from all sessions", "user", user.ID)
	return nil
}
