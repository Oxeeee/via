package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OxytocinGroup/theca-v3/internal/model"
	customerrors "github.com/OxytocinGroup/theca-v3/internal/utils/errors"
	"github.com/OxytocinGroup/theca-v3/internal/vars"
	"gorm.io/gorm"
)

type Repository interface {
	Register(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
}

type repository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewRepository(db *gorm.DB, log *slog.Logger) Repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (r *repository) Register(ctx context.Context, user *model.User) error {
	const op = "repository.Register"
	log := r.log.With("op", op)

	var c int64
	err := r.db.Model(&model.User{}).Where("username = ?", user.Username).Count(&c).Error
	if err != nil {
		log.Error("failed to count users", "error", err)
		return customerrors.FromGormError(err)
	}

	if c > 0 {
		return customerrors.New(customerrors.CodeUserAlreadyExists, "Пользователь с таким именем уже существует")
	}

	err = r.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&c).Error
	if err != nil {
		log.Error("failed to count users", "error", err)
		return customerrors.FromGormError(err)
	}
	if c > 0 {
		return customerrors.New(customerrors.CodeUserAlreadyExists, "Пользователь с такой почтой уже существует")
	}

	err = r.db.Model(&model.User{}).Create(user).Error
	if err != nil {
		log.Error("failed to create user", "error", err)
		return customerrors.FromGormError(err)
	}

	return nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	const op = "repository.GetUserByUsername"
	log := r.log.With("op", op)

	var user model.User
	err := r.db.Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.New(customerrors.CodeUserNotFound, "Пользователь не найден")
		}
		log.Error("failed to get user by username", "error", err)
		return nil, customerrors.FromGormError(err)
	}

	return &user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	const op = "repository.GetUserByID"
	log := r.log.With("op", op)

	var user model.User
	err := r.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, vars.ErrUserNotFound
		}
		log.Error("failed to get user by id", "error", err)
		return nil, err
	}
	return &user, nil
}

func (r *repository) SaveUser(ctx context.Context, user *model.User) error {
	const op = "repository.SaveUser"
	log := r.log.With("op", op)

	err := r.db.Model(&model.User{}).Save(user).Error
	if err != nil {
		log.Error("failed to save user", "error", err)
		return err
	}

	return nil
}
